package callbacks

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func RegisterDefaultCallbacks(db *gorm.DB, config *callbacks.Config) {
	createCallback := db.Callback().Create()
	createCallback.Register("arango:create", Create(config))

	queryCallback := db.Callback().Query()
	queryCallback.Register("arango:query", Query)

	updateCallback := db.Callback().Update()
	updateCallback.Register("arango:update", Update)
}
