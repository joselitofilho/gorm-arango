package conn

import (
	"context"
	"database/sql"
	"errors"

	driver "github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm/driver/arango/internal/aql"
	"github.com/joselitofilho/gorm/driver/arango/internal/transformers"
)

type ConnPool struct {
	Connection driver.Connection
	Database   driver.Database
}

func (connPool *ConnPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

func (connPool *ConnPool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	vars := make(map[string]interface{}, len(args)/2)

	arangoResult := args[len(args)-1]

	transformers.SliceToMap(args[:len(args)-1], vars)
	cursor, err := connPool.Database.Query(ctx, query, vars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	// ctx = driver.WithQueryFullCount(ctx)
	_, err = cursor.ReadDocument(ctx, arangoResult)
	if driver.IsNoMoreDocuments(err) {
		return nil, errors.New("Document not found")
	} else if err != nil {
		return nil, err
	}

	result := aql.NewResult(0, cursor.Count(), arangoResult, nil)
	return result, nil
}

func (connPool *ConnPool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (connPool *ConnPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}
