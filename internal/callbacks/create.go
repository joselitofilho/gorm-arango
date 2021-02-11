package callbacks

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"

	"github.com/joselitofilho/gorm/driver/arango/internal/conn"
)

func Create(config *callbacks.Config) func(db *gorm.DB) {
	if config.WithReturning {
		return callbacks.CreateWithReturning
	} else {
		return func(db *gorm.DB) {
			if db.Error == nil {
				connection := db.Statement.ConnPool.(*conn.ConnPool)

				collection, err := connection.Database.Collection(db.Statement.Context, db.Statement.Table)
				if err == nil {
					entity := map[string]interface{}{}
					data, _ := json.Marshal(db.Statement.Dest)
					_ = json.Unmarshal(data, &entity)
					// TODO: generate ID and meta
					_, err = collection.CreateDocument(db.Statement.Context, entity)
					if err == nil {
						db.RowsAffected = int64(1)
					} else {
						db.AddError(err)
					}
				} else {
					db.AddError(err)
				}
			}
		}
	}
}
