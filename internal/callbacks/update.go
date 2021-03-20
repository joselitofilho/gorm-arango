package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm/driver/arango/internal/conn"
	"gorm.io/gorm"
)

func Update(db *gorm.DB) {
	if !db.DryRun && db.Error == nil {
		connection := db.Statement.ConnPool.(*conn.ConnPool)
		if collection, err := connection.Database.Collection(db.Statement.Context, db.Statement.Table); err == nil {
			entityMap := map[string]interface{}{}
			data, _ := json.Marshal(db.Statement.Dest)
			_ = json.Unmarshal(data, &entityMap)

			bindKeysVars := map[string]interface{}{"ID": entityMap["ID"]}
			docMetaInfo, err := getMeta(db.Statement.Context, collection, bindKeysVars, db.Statement.Dest)
			if err != nil {
				db.AddError(err)
				return
			}

			entityMap["UpdatedAt"] = time.Now()

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
	filters := " FILTER doc.DeletedAt == null "
	for key := range bindKeyVars {
		filters += fmt.Sprintf(" AND doc.%s == @%s ", key, key)
	}

	query := fmt.Sprintf("FOR doc IN %s %s RETURN doc", collection.Name(), filters)

	cursor, err := collection.Database().Query(ctx, query, bindKeyVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	documentMetaInfo, err := cursor.ReadDocument(ctx, result)
	if driver.IsNoMoreDocuments(err) {
		// TODO: Create a better error
		return nil, fmt.Errorf("Entity not found")
	} else if err != nil {
		return nil, err
	}

	return &documentMetaInfo, nil
}