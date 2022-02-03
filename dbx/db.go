package dbx

import (
	"context"
	"database/sql"
	"dcamachoj/time-tracker-rest/common"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var dbCurr *sql.DB
var dbType = common.GetEnv("DB_TYPE", "mysql")
var dbUses = common.GetEnv("DB_USES", "db-test")
var dbUser = common.GetEnv("DB_USER", "db-user")
var dbPass = common.GetEnv("DB_PASS", "db-pass")
var dbHost = common.GetEnv("DB_HOST", "localhost")
var dbPort = common.GetEnv("DB_PORT", "3306")
var dbArgs = common.GetEnv("DB_ARGS", "?charset=utf8mb4&parseTime=True&loc=Local")
var dbConn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", dbUser, dbPass, dbHost, dbPort, dbUses, dbArgs)

func GetSchema() string {
	return dbUses
}

func getDb() (*sql.DB, error) {
	var err error
	var db = dbCurr
	if db != nil {
		err = db.Ping()
		if err != nil {
			return db, nil
		}
		db.Close()
	}
	// fmt.Println("DB Connection ", dbConn)
	db, err = sql.Open(dbType, dbConn)
	if err != nil {
		return nil, errors.Wrapf(err, "Opening DB %s: %s", dbType, dbConn)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "Pinging DB %s: %s", dbType, dbConn)
	}
	dbCurr = db
	return db, nil
}

func JoinSQL(lines ...string) string {
	return strings.Join(lines, "\r\n")
}
func ExecSQLs(ctx context.Context, sqlStrs ...string) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var tx = GetTx(ctx)
		for k, s := range sqlStrs {
			if s == "" {
				continue
			}
			_, err = tx.Exec(s)
			if err != nil {
				return errors.Wrapf(err, "Executing SQL %d: %s", k, s)
			}
		}
		return err
	})

	return err
}
func ExecTx(ctx context.Context, handler TxHandlerFunc) error {
	var err error
	var db *sql.DB
	var tx *sql.Tx

	if GetTx(ctx) != nil {
		// already inside db context
		return handler(ctx)
	}

	db, err = getDb()
	if err != nil {
		return errors.Wrap(err, "Get DB")
	}

	tx, err = db.Begin()
	if err != nil {
		return errors.Wrap(err, "Begin")
	}
	defer tx.Rollback()

	err = handleTx(ctx, handler, tx)
	if err != nil {
		if re, ok := err.(*common.Response); ok {
			// if error is a response to the user, don't wrap
			return re
		}
		return errors.Wrap(err, "Handling Tx")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "Commit")
	}

	return nil
}

func handleTx(ctx context.Context, handler TxHandlerFunc, tx *sql.Tx) (err error) {
	defer func() {
		if rErr := common.RecoverError(recover(), "handleTx"); rErr != nil {
			err = rErr
		}
	}()
	return handler(WithTx(ctx, tx))
}
func LikeString(str string, startStar, endStar bool) string {
	if str == "" {
		if startStar || endStar {
			return "%"
		}
		return ""
	}
	if startStar {
		str = "*" + str
	}
	if endStar {
		str += "*"
	}
	str = strings.Replace(str, "*", "%", -1)
	str = strings.Replace(str, " ", "%", -1)
	str = strings.Replace(str, "?", "_", -1)
	for strings.Contains(str, "**") {
		str = strings.Replace(str, "**", "*", -1)
	}
	return str
}
func SqlPage(query string, limit int, offset int) string {
	return fmt.Sprintf("%s\r\nLIMIT %d OFFSET %d", query, limit, offset)
}
func dbLog(query string, args ...interface{}) {
	const indent = "  "
	const crlf = "\r\n"
	var sb strings.Builder
	sb.WriteString("Query:" + crlf)
	sb.WriteString(query + crlf)
	sb.WriteString("Args:" + crlf)
	for _, arg := range args {
		sb.WriteString(indent + fmt.Sprint(arg) + crlf)
	}
	log.Println(sb.String())
}

func GetScalar(ctx context.Context, dest interface{}, isNull *bool,
	query string, args ...interface{}) error {
	var tx = GetTx(ctx)
	if dest == nil {
		return errors.Errorf("dest can't be nil")
	}
	var pDest = reflect.ValueOf(dest)
	var err = checkPtrScalar(pDest.Type())
	if err != nil {
		return err
	}

	var rows *sql.Rows

	rows, err = tx.Query(query, args...)
	if err != nil {
		err = errors.Wrap(err, "ScalarGet")
		dbLog(query, args...)
		log.Println(err)
		return err
	}
	defer rows.Close()

	if rows.Next() {

		err = rows.Scan(dest)
		if err != nil {
			return err
		}

	} else {
		if isNull != nil {
			*isNull = true
		}
		return nil
	}

	if isNull != nil {
		*isNull = false
	}
	return nil
}
func ListScalar(ctx context.Context, pList interface{},
	query string, args ...interface{}) error {
	var tx = GetTx(ctx)
	if pList == nil {
		return errors.Errorf("dest can't be nil")
	}
	var pDest = reflect.ValueOf(pList)
	var err = checkPtrSliceScalar(pDest.Type())
	if err != nil {
		return err
	}

	var sDest = pDest.Elem()
	var tDest = pDest.Type().Elem().Elem()
	var row interface{}
	var rows *sql.Rows

	rows, err = tx.Query(query, args...)
	if err != nil {
		err = errors.Wrap(err, "ScalarList")
		dbLog(query, args...)
		log.Println(err)
		return err
	}
	defer rows.Close()

	var index = 0
	for rows.Next() {
		row = reflect.New(tDest).Interface()
		err = rows.Scan(row)
		if err != nil {
			return errors.Wrapf(err, "Scanning Row %d", index)
		}

		index++
		sDest.Set(reflect.Append(sDest, reflect.ValueOf(row).Elem()))
	}

	return nil
}

func PageScalar(ctx context.Context, pList interface{},
	page *common.ResponsePage,
	queryLst, queryCnt string,
	args ...interface{}) error {
	var err error
	if page.Size == 0 {
		page.Size = common.DefaultPage
	}
	queryLst = SqlPage(queryLst, page.Size, page.Size*page.Index)

	var cnt int64
	var isNull bool

	err = GetScalar(ctx, &cnt, &isNull, queryCnt, args...)
	if err != nil {
		return errors.Wrap(err, "Page.Count")
	}
	if isNull {
		cnt = 0
	}

	err = ListScalar(ctx, pList, queryLst, args...)
	if err != nil {
		return errors.Wrap(err, "Page.List")
	}

	page.SetTotal(cnt)
	return nil
}
func GetEntity(ctx context.Context,
	entity Entity, dest interface{},
	query string, args ...interface{}) error {
	if dest == nil {
		return errors.Errorf("dest can't be nil")
	}
	var ppDest = reflect.ValueOf(dest)
	var err = checkType(ppDest.Type(), "a pointer to an Entity pointer", reflect.Ptr, reflect.Ptr)
	if err != nil {
		return err
	}
	var rows *sql.Rows
	var scanner = newScanner(query, args...)
	var row interface{}
	var pDest = ppDest.Elem()

	var tx = GetTx(ctx)
	rows, err = tx.Query(query, args...)
	if err != nil {
		err = errors.Wrap(err, "EntityGet")
		dbLog(query, args...)
		log.Println(err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		row, err = scanner.Scan(ctx, rows, entity)
		if err != nil {
			return err
		}

		if row == nil {
			pDest.Set(reflect.Zero(pDest.Type()))
			return nil
		}

		pDest.Set(reflect.ValueOf(row))

	} else {
		pDest.Set(reflect.Zero(pDest.Type()))
		return nil
	}

	return nil
}

func ListEntity(ctx context.Context,
	entity Entity,
	pList interface{},
	query string, args ...interface{}) error {
	if pList == nil {
		return errors.Errorf("dest can't be nil")
	}

	var pDest = reflect.ValueOf(pList)
	if !isPointer(pDest.Type()) {
		return errors.Errorf("dest should be a pointer to a list of Entity Pointers. Actual: %T", pList)
	}
	if !isSlice(pDest.Type().Elem()) {
		return errors.Errorf("dest should be a pointer to a list of Entity Pointers. Actual: %T", pList)
	}
	if !isPointer(pDest.Type().Elem().Elem()) {
		return errors.Errorf("dest should be a pointer to a list of Entity Pointers. Actual: %T", pList)
	}

	var err error
	var sDest = pDest.Elem()
	var scanner = newScanner(query, args...)
	var row interface{}
	var rows *sql.Rows

	var tx = GetTx(ctx)
	rows, err = tx.Query(query, args...)
	if err != nil {
		err = errors.Wrap(err, "EntityList")
		dbLog(query, args...)
		log.Println(err)
		return err
	}
	defer rows.Close()

	var index = 0
	for rows.Next() {
		row, err = scanner.Scan(ctx, rows, entity)
		if err != nil {
			return errors.Wrapf(err, "Scanning Row %d", index)
		}

		index++
		sDest.Set(reflect.Append(sDest, reflect.ValueOf(row)))
	}

	return nil
}
func PageEntity(ctx context.Context,
	entity Entity, pList interface{}, page *common.ResponsePage,
	queryLst, queryCnt string, args ...interface{}) error {
	var err error
	if page.Size == 0 {
		page.Size = common.DefaultPage
	}
	queryLst = SqlPage(queryLst, page.Size, page.Size*page.Index)

	var cnt int64
	var isNull bool

	err = GetScalar(ctx, &cnt, &isNull, queryCnt, args...)
	if err != nil {
		return errors.Wrap(err, "Page.Count")
	}
	if isNull {
		cnt = 0
	}

	err = ListEntity(ctx, entity, pList, queryLst, args...)
	if err != nil {
		return errors.Wrap(err, "Page.List")
	}

	page.SetTotal(cnt)
	return nil
}
