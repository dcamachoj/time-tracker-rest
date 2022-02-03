package dbx

import (
	"context"
	"database/sql"
	"dcamachoj/time-tracker-rest/common"
	"strings"
)

// Employee DTO
type Employee struct {
	ID          int64       `json:"id,omitempty"`
	UserName    string      `json:"user_name,omitempty"`
	DisplayName string      `json:"display_name,omitempty"`
	TimeData    []*TimeData `json:"time_data,omitempty"`
}

type EmployeeEntity struct {
	ID          sql.NullInt64
	UserName    sql.NullString
	DisplayName sql.NullString
}

var employeeEntity = &EmployeeEntity{}

func (e *EmployeeEntity) PreScan(scanner EntityScanner) error {
	scanner.AddScanners(
		&e.ID,
		&e.UserName,
		&e.DisplayName,
	)
	return nil
}
func (e *EmployeeEntity) PosScan(scanner EntityScanner) error {
	return nil
}
func (e *EmployeeEntity) Result() interface{} {
	if !e.ID.Valid {
		return nil
	}
	return &Employee{
		ID:          e.ID.Int64,
		UserName:    e.UserName.String,
		DisplayName: e.DisplayName.String,
	}
}
func (e *EmployeeEntity) Tablename() string {
	return "employees"
}
func (e *EmployeeEntity) Select(prefix string) []string {
	if prefix == "" {
		prefix = "tb."
	}
	return []string{
		prefix + "id",
		prefix + "user_name",
		prefix + "display_name",
	}
}
func (e *EmployeeEntity) Page(ctx context.Context, pList *[]*Employee, page *common.ResponsePage) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var qryLst = NewBuilder().
			Addln("SELECT "+strings.Join(e.Select("tb."), ", ")).
			Addf("FROM %s AS tb", e.Tablename()).
			Addln().
			Addln("ORDER BY tb.user_name").
			String()
		var qryCnt = NewBuilder().
			Addln("SELECT COUNT(tb.id)").
			Addf("FROM %s AS tb", e.Tablename()).
			Addln().
			String()
		err = PageEntity(ctx, e, pList, page, qryLst, qryCnt)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
func (e *EmployeeEntity) GetSingle(ctx context.Context, id int64, pRes **Employee) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var qry = NewBuilder().Addln("SELECT "+strings.Join(e.Select("tb."), ", ")).
			Addf("FROM %s AS tb", e.Tablename()).Addln().
			Addln("WHERE tb.id = ?").
			String()
		err = GetEntity(ctx, e, pRes, qry, id)
		if err != nil {
			return nil
		}
		if *pRes == nil { // not found
			return nil
		}

		var timePage = &common.ResponsePage{
			Size: 10,
		}
		var timeData = (*pRes).TimeData
		err = timeDataEntity.Page(ctx, id, &timeData, timePage)
		if err != nil {
			return nil
		}

		return nil
	})
	return err
}
func (e *EmployeeEntity) Insert(ctx context.Context, row *Employee) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var tx = GetTx(ctx)
		var qry = NewBuilder().
			Addf("INSERT INTO %s (user_name, display_name)", e.Tablename()).Addln().
			Addln("VALUES (?, ?);").
			String()
		var res sql.Result
		res, err = tx.Exec(qry, row.UserName, row.DisplayName)
		if err != nil {
			return err
		}
		row.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeEntity) Update(ctx context.Context, row *Employee) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var tx = GetTx(ctx)
		var qry = NewBuilder().
			Addf("UPDATE %s", e.Tablename()).Addln().
			Addln("SET display_name = ?").
			Addln("WHERE id = ?").
			String()
		_, err = tx.Exec(qry, row.DisplayName, row.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeEntity) Delete(ctx context.Context, id int64) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var tx = GetTx(ctx)
		var qry = NewBuilder().
			Addf("DELETE FROM %s", e.Tablename()).Addln().
			Addln("WHERE id = ?").
			String()
		_, err = tx.Exec(qry, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
