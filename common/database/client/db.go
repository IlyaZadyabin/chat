package client

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler - function that executes within a transaction
type Handler func(ctx context.Context) error

// Client interface for working with database
type Client interface {
	DB() DB
	Close() error
}

// TxManager transaction manager that executes user-specified handler within a transaction
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query wrapper over a query, storing the query name and the query itself
// Query name is used for logging and potentially can be used elsewhere, for example, for tracing
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor interface for working with transactions
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer combines NamedExecer and QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer interface for working with named queries using struct tags
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer interface for working with regular queries
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger interface for checking database connection
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB interface for working with database
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
