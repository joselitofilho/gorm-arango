package arango

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

// Migrator ...
type Migrator struct {
	migrator.Migrator
}

// CurrentDatabase ...
func (m Migrator) CurrentDatabase() (name string) {
	if dialector, ok := m.DB.Dialector.(Dialector); ok {
		name = dialector.Database.Name()
	}
	return
}

// HasTable ...
func (m Migrator) HasTable(value interface{}) bool {
	var hasTable bool
	var err error

	err = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if dialector, ok := m.DB.Dialector.(Dialector); ok {
			currentDatabase := m.DB.Migrator().CurrentDatabase()
			fmt.Println(currentDatabase)
			hasTable, err = dialector.CollectionExists(stmt.Table)
			return err
		}
		return ErrDatabaseConnectionFailed
	})
	if err != nil {
		panic(err)
	}

	return hasTable
}

// CreateTable ...
func (m Migrator) CreateTable(values ...interface{}) error {
	return m.RunWithValue(values[0], func(stmt *gorm.Statement) error {
		if dialector, ok := m.DB.Dialector.(Dialector); ok {
			_, err := dialector.CreateCollection(stmt.Table)
			return err
		}
		return ErrDatabaseConnectionFailed
	})
}
