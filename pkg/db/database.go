package db

import (
	"context"
	"database/sql"
)

// Connection - provides functions to interact with an underlying database
type Connection interface {
	Open() error
	Close() error

	PrepareContext(ctx context.Context, query string) (Execable, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Scannable
	QueryContext(ctx context.Context, query string, args ...interface{}) (ScannableList, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Execable - provides an interface for a statement to be executed
type Execable interface {
	Exec(args ...interface{}) (sql.Result, error)
}

// Scannable - provides an interface for results can be parsed
type Scannable interface {
	Scan(dest ...interface{}) error
}

// ScannableList - provides an interfaces for a list of results
type ScannableList interface {
	Scannable
	Close() error
	Next() bool
	Err() error
}

// NewConnection -
func NewConnection(adapter, dsn string) Connection {
	return &connection{
		adapter: adapter,
		dsn:     dsn,
	}
}

type connection struct {
	adapter string
	dsn     string
	session *sql.DB
}

func (c *connection) Open() error {
	sess, err := sql.Open(c.adapter, c.dsn)
	if err != nil {
		return err
	}

	c.session = sess
	return nil
}

func (c *connection) Close() error {
	return c.session.Close()
}

func (c *connection) connection() *sql.DB {
	return c.session
}

func (c *connection) PrepareContext(ctx context.Context, query string) (Execable, error) {
	return c.connection().PrepareContext(ctx, query)
}

func (c *connection) QueryRowContext(ctx context.Context, query string, args ...interface{}) Scannable {
	return c.connection().QueryRowContext(ctx, query, args...)
}

func (c *connection) QueryContext(ctx context.Context, query string, args ...interface{}) (ScannableList, error) {
	return c.connection().QueryContext(ctx, query, args...)
}
func (c *connection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.connection().ExecContext(ctx, query, args...)
}
