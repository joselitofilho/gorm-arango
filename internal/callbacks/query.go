package callbacks

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm/driver/arango/internal/conn"
	"gorm.io/gorm"
)

func Query(db *gorm.DB) {
	if db.Error == nil {
		connection := db.Statement.ConnPool.(*conn.ConnPool)
		if collection, err := connection.Database.Collection(db.Statement.Context, db.Statement.Table); err == nil {
			query, vars, err := BuildAQL(db)
			if err != nil {
				db.AddError(err)
				return
			}

			cursor, err := collection.Database().Query(db.Statement.Context, query, vars)
			if err != nil {
				db.AddError(err)
				return
			}
			defer cursor.Close()
			_, err = cursor.ReadDocument(db.Statement.Context, db.Statement.Dest)
			if driver.IsNoMoreDocuments(err) {
				db.AddError(errors.New("Document not found"))
				return
			} else if err != nil {
				db.AddError(err)
				return
			}
		}
	}
}

// BuildAQL ...
func BuildAQL(db *gorm.DB) (query string, bindingFields map[string]interface{}, err error) {
	bindingFields = map[string]interface{}{"@collection": db.Statement.Table}

	db.Statement.Build("OFFSET")
	offsetStr := strings.TrimSpace(strings.ReplaceAll(db.Statement.SQL.String(), "OFFSET", ""))
	db.Statement.SQL.Reset()
	offset := 0
	if len(offsetStr) > 0 {
		if offset, err = strconv.Atoi(offsetStr); err != nil {
			offset = 0
		}
	}
	bindingFields["offset"] = offset

	db.Statement.Build("LIMIT")
	limitStr := strings.TrimSpace(strings.ReplaceAll(db.Statement.SQL.String(), "LIMIT", ""))
	db.Statement.SQL.Reset()
	limit := 250
	if len(limitStr) > 0 {
		if limit, err = strconv.Atoi(limitStr); err != nil {
			limit = 250
		}
	}
	bindingFields["limit"] = limit

	db.Statement.Build("WHERE")
	filters := db.Statement.SQL.String()

	// TODO: To be continued...
	query = fmt.Sprintf("FOR doc IN @@collection FILTER doc.DeleteAt == null %s LIMIT @offset, @limit RETURN doc", filters)

	// TODO: sort.

	return
}
