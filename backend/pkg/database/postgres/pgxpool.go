package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PgxPool is a PostgreSQL connection pool.
// abstracts the PostgreSQL connection pool. It contains methods for executing queries and managing the pool.
type PgxPool interface { //we heed ti for implementating postgre interface
	Close()
	// Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row

	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}
