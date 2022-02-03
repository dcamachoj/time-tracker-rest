package dbx

import (
	"context"
	"database/sql"
)

const ckTx = ckDbx("tx")

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ckTx, tx)
}

func GetTx(ctx context.Context) *sql.Tx {
	var tx, ok = ctx.Value(ckTx).(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}
