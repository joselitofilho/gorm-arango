package callbacks

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"

	"github.com/joselitofilho/gorm/driver/arango/internal/conn"
	"github.com/mitchellh/mapstructure"
)

func Create(config *callbacks.Config) func(db *gorm.DB) {
	if config.WithReturning {
		return callbacks.CreateWithReturning
	} else {
		return func(db *gorm.DB) {
			if db.Error == nil {
				connection := db.Statement.ConnPool.(*conn.ConnPool)
				if collection, err := connection.Database.Collection(db.Statement.Context, db.Statement.Table); err == nil {
					entityMap := map[string]interface{}{}
					data, _ := json.Marshal(db.Statement.Dest)
					_ = json.Unmarshal(data, &entityMap)

					modelMap := map[string]interface{}{}
					now := time.Now()
					model := gorm.Model{
						ID:        uint(now.UnixNano()), // TODO: check this later...
						CreatedAt: now,
						UpdatedAt: now,
					}
					data, _ = json.Marshal(model)
					_ = json.Unmarshal(data, &modelMap)

					for k, v := range modelMap {
						entityMap[k] = v
					}

					if _, err = collection.CreateDocument(db.Statement.Context, entityMap); err == nil {
						entityMap["Model"] = model
						mapstructure.Decode(entityMap, &db.Statement.Dest)
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
