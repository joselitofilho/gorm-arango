package conn

import (
	"context"
	"database/sql"
	"errors"

	driver "github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm/driver/arango/internal/transformers"
)

type ConnPool struct {
	Connection driver.Connection
	Database   driver.Database
	Result     interface{}
}

func (connPool *ConnPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// TODO: Implements
	return nil, nil
}

func (connPool *ConnPool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// TODO: Implements
	return nil, nil
}

func (connPool *ConnPool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	vars := transformers.SliceToMap(args)
	cursor, err := connPool.Database.Query(ctx, query, vars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	_, err = cursor.ReadDocument(ctx, connPool.Result)

	if driver.IsNoMoreDocuments(err) {
		return nil, errors.New("Document not found")
	} else if err != nil {
		return nil, err
	}

	// TODO: need implements parse to *sql.Rows?
	return nil, nil
}

func (connPool *ConnPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// TODO: Implements
	return nil
}
