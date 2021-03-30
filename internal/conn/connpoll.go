package conn

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	driver "github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm-arango/internal/transformers"
)

type ConnPoolReturn struct {
	Dest     interface{}
	ElemType reflect.Type
	IsSlice  bool
}

type ConnPool struct {
	Connection driver.Connection
	Database   driver.Database
	Return     ConnPoolReturn
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

	if connPool.Return.IsSlice {
		results := make([]interface{}, 0)
		for {
			var r interface{}
			r = reflect.New(connPool.Return.ElemType).Interface()
			_, err = cursor.ReadDocument(ctx, r)

			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				return nil, err
			} else {
				results = append(results, r)
			}
		}
		connPool.Return.Dest = results
	} else {
		_, err = cursor.ReadDocument(ctx, connPool.Return.Dest)

		if driver.IsNoMoreDocuments(err) {
			return nil, errors.New("Document not found")
		} else if err != nil {
			return nil, err
		}
	}

	// TODO: need implements parse to *sql.Rows?
	return nil, nil
}

func (connPool *ConnPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// TODO: Implements
	return nil
}
