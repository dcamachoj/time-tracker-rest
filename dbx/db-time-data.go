package dbx

import (
	"context"
	"database/sql"
	"dcamachoj/time-tracker-rest/common"
	"strings"
)

// TimeData
type TimeData struct {
	ID         int64      `json:"id,omitempty"`
	ParentID   int64      `json:"parent_id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Title      string     `json:"title,omitempty"`
	TimeFrames TimeFrames `json:"frames,omitempty"`
}

// TimeDataEntity
type TimeDataEntity struct {
	ID         sql.NullInt64
	ParentID   sql.NullInt64
	Name       sql.NullString
	Title      sql.NullString
	TimeFrames TimeFramesEntity
}

var timeDataEntity = &TimeDataEntity{}

func (e *TimeDataEntity) PreScan(scanner EntityScanner) error {
	scanner.AddScanners(
		&e.ID,
		&e.ParentID,
		&e.Name,
		&e.Title,
	)
	e.TimeFrames.PreScan(scanner)
	return nil
}
func (e *TimeDataEntity) PosScan(scanner EntityScanner) error {
	return nil
}
func (e *TimeDataEntity) Result() interface{} {
	if !e.ID.Valid {
		return nil
	}
	var res = &TimeData{
		ID:       e.ID.Int64,
		ParentID: e.ParentID.Int64,
		Name:     e.Name.String,
		Title:    e.Title.String,
	}
	e.TimeFrames.AsTimeFrames(&res.TimeFrames)
	return res
}
func (e *TimeDataEntity) Tablename() string {
	return "time_data"
}
func (e *TimeDataEntity) Select(prefix string) []string {
	var res = []string{
		prefix + "id",
		prefix + "parent_id",
		prefix + "name",
		prefix + "title",
	}
	res = append(res, e.TimeFrames.Select(prefix)...)

	return res
}
func (e *TimeDataEntity) Page(ctx context.Context,
	parentId int64,
	pList *[]*TimeData, page *common.ResponsePage) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var qryLst = NewBuilder().
			Addln("SELECT "+strings.Join(e.Select("tb."), ", ")).
			Addf("FROM %s AS tb", e.Tablename()).Addln().
			Addln("WHERE tb.parent_id = ?").
			Addln("ORDER BY tb.name").
			String()
		var qryCnt = NewBuilder().
			Addln("SELECT COUNT(tb.id)").
			Addf("FROM %s AS tb", e.Tablename()).Addln().
			Addln("WHERE tb.parent_id = ?").
			String()
		err = PageEntity(ctx, e, pList, page, qryLst, qryCnt, parentId)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (e *TimeDataEntity) GetSingle(ctx context.Context, id int64, pRes **TimeData) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		// load TimeData
		var qryData = NewBuilder().
			Addln("SELECT "+strings.Join(e.Select("tb."), ", ")).
			Addf("FROM %s AS tb", e.Tablename()).Addln().
			Addln("WHERE tb.id = ?").
			String()
		err = GetEntity(ctx, e, pRes, qryData, id)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (e *TimeDataEntity) Insert(ctx context.Context, row *TimeData) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var tx = GetTx(ctx)
		var qry = NewBuilder().
			Addf("INSERT INTO %s (%s)", e.Tablename(), strings.Join(e.Select("")[1:], ", ")).Addln().
			Addln("VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);").
			String()
		var res sql.Result
		res, err = tx.Exec(qry,
			row.ParentID,
			row.Name,
			row.Title,
			row.TimeFrames.Daily.Current,
			row.TimeFrames.Daily.Previous,
			row.TimeFrames.Weekly.Current,
			row.TimeFrames.Weekly.Previous,
			row.TimeFrames.Monthly.Current,
			row.TimeFrames.Monthly.Previous,
		)
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

func (e *TimeDataEntity) Update(ctx context.Context, row *TimeData) error {
	var err error
	err = ExecTx(ctx, func(ctx context.Context) error {
		var tx = GetTx(ctx)
		var qb = NewBuilder().
			Addf("UPDATE %s", e.Tablename()).Addln()
		qb.Addln("SET")
		var fieldNames = e.Select("")[4:]
		var sep = ","
		for k, name := range fieldNames {
			if k == 5 {
				sep = ""
			}
			qb.Addf("%s  = ?%s", name, sep).Addln()
		}
		var qry = qb.
			Addf("WHERE id = ?").Addln().
			String()
		_, err = tx.Exec(qry,
			row.TimeFrames.Daily.Current,
			row.TimeFrames.Daily.Previous,
			row.TimeFrames.Weekly.Current,
			row.TimeFrames.Weekly.Previous,
			row.TimeFrames.Monthly.Current,
			row.TimeFrames.Monthly.Previous,
			row.ID,
		)
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

func (e *TimeDataEntity) Delete(ctx context.Context, id int64) error {
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
