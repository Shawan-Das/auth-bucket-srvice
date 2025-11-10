package db

import (
	"github.com/jackc/pgx/v4"
)

type DBTX interface {
	Exec(ctx interface{}, sql string, arguments ...interface{}) (pgx.CommandTag, error)
	Query(ctx interface{}, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx interface{}, sql string, args ...interface{}) pgx.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx DBTX) *Queries {
	return &Queries{
		db: tx,
	}
}
