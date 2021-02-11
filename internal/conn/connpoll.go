package conn

import (
	"context"
	"database/sql"

	driver "github.com/arangodb/go-driver"
)

type ConnPool struct {
	Connection driver.Connection
	Database   driver.Database
}

func (connPool *ConnPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

func (connPool *ConnPool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (connPool *ConnPool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (connPool *ConnPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}
