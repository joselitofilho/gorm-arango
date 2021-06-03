package callbacks

import (
	"context"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm-arango/internal/conn"
	"github.com/joselitofilho/gorm-arango/internal/transformers"
	"gorm.io/gorm"
)

func Update(db *gorm.DB) {
	if !db.DryRun && db.Error == nil {
		connection := db.Statement.ConnPool.(*conn.ConnPool)
		if collection, err := connection.Database.Collection(db.Statement.Context, db.Statement.Table); err == nil {
			modelMap, err := transformers.EntityToMap(db.Statement.Model)
			if err != nil {
				db.AddError(err)
				return
			}

			entityMap, err := transformers.EntityToMap(db.Statement.Dest)
			if err != nil {
				db.AddError(err)
				return
			}

			bindKeysVars := map[string]interface{}{"ID": modelMap["ID"]}
			docMetaInfo, err := getMeta(db.Statement.Context, collection, bindKeysVars, db.Statement.Model)
			if err != nil {
				db.AddError(err)
				return
			}

			entityMap["UpdatedAt"] = time.Now()
			delete(entityMap, "ID")
			delete(entityMap, "CreatedAt")
			delete(entityMap, "DeletedAt")

			if _, err := collection.UpdateDocument(db.Statement.Context, docMetaInfo.Key, entityMap); err != nil {
				db.AddError(err)
			}

			db.RowsAffected = 1

			// result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			// if err == nil {
			// 	db.RowsAffected, _ = result.RowsAffected()
			// } else {
			// 	db.AddError(err)
			// }
		} else {
			db.AddError(err)
		}
	}
}

func getMeta(ctx context.Context, collection driver.Collection, bindKeyVars map[string]interface{}, result interface{}) (*driver.DocumentMeta, error) {
	aliasWithoutDot := collection.Name()
	alias := aliasWithoutDot + "."

	filters := fmt.Sprintf(" FILTER %sDeletedAt == null ", alias)
	for key := range bindKeyVars {
		filters += fmt.Sprintf(" AND %s%s == @%s ", alias, key, key)
	}

	query := fmt.Sprintf("FOR %s IN %s %s RETURN %s", aliasWithoutDot, aliasWithoutDot, filters, aliasWithoutDot)

	cursor, err := collection.Database().Query(ctx, query, bindKeyVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	documentMetaInfo, err := cursor.ReadDocument(ctx, result)
	if driver.IsNoMoreDocuments(err) {
		// TODO: Create a better error
		return nil, fmt.Errorf("entity not found")
	} else if err != nil {
		return nil, err
	}

	return &documentMetaInfo, nil
}
