package callbacks

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

func RegisterDefaultCallbacks(db *gorm.DB, config *callbacks.Config) {
	createCallback := db.Callback().Create()
	createCallback.Register("gorm:create", Create(config))
}
