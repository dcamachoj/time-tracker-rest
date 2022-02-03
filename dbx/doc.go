package dbx

import (
	"context"
	"database/sql"
	"dcamachoj/time-tracker-rest/common"
	"time"
)

type ckDbx string
type TxHandlerFunc func(ctx context.Context) error

// Entity interface
type Entity interface {
	PreScan(scanner EntityScanner) error
	PosScan(scanner EntityScanner) error
	Result() interface{}
}

// EntityTable interface
type EntityTable interface {
	Tablename() string
}

// EntitySelect interface
type EntitySelect interface {
	Select(prefix string) []string
}

// EntityCtor type
type EntityCtor func() Entity

// EntityScanner interface
type EntityScanner interface {
	Scanners() []interface{}
	AddScanners(scanners ...interface{})
	Scan(ctx context.Context, rows *sql.Rows, entity Entity) (interface{}, error)
}

// SqlString pointer
func SqlString(sn sql.NullString) *string {
	if !sn.Valid {
		return nil
	}

	return common.PtrString(sn.String)
}

// SqlBool pointer
func SqlBool(sn sql.NullBool) *bool {
	if !sn.Valid {
		return nil
	}

	return common.PtrBool(sn.Bool)
}

// SqlInt32 pointer
func SqlInt32(sn sql.NullInt32) *int32 {
	if !sn.Valid {
		return nil
	}

	return common.PtrInt32(sn.Int32)
}

// SqlInt64 pointer
func SqlInt64(sn sql.NullInt64) *int64 {
	if !sn.Valid {
		return nil
	}

	return common.PtrInt64(sn.Int64)
}

// SqlFloat64 pointer
func SqlFloat64(sn sql.NullFloat64) *float64 {
	if !sn.Valid {
		return nil
	}

	return common.PtrFloat64(sn.Float64)
}

// SqlTime pointer
func SqlTime(sn sql.NullTime) *time.Time {
	if !sn.Valid {
		return nil
	}

	return common.PtrTime(sn.Time)
}
