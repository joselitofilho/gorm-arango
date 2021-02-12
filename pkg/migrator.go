package arango

import (
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

// Migrator ...
type Migrator struct {
	migrator.Migrator
}

// AutoMigrate ...
func (m Migrator) AutoMigrate(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, true) {
		tx := m.DB.Session(&gorm.Session{})
		if !tx.Migrator().HasTable(value) {
			if err := tx.Migrator().CreateTable(value); err != nil {
				return err
			}
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Database
////////////////////////////////////////////////////////////////////////////////

// CurrentDatabase ...
func (m Migrator) CurrentDatabase() (name string) {
	if dialector, ok := m.DB.Dialector.(Dialector); ok {
		name = dialector.Database.Name()
	}
	return
}

////////////////////////////////////////////////////////////////////////////////
// Tables
////////////////////////////////////////////////////////////////////////////////

// CreateTable ...
func (m Migrator) CreateTable(values ...interface{}) error {
	return m.RunWithValue(values[0], func(stmt *gorm.Statement) error {
		if dialector, ok := m.DB.Dialector.(Dialector); ok {
			_, err := dialector.CreateCollection(stmt.Context, stmt.Table)
			return err
		}
		return ErrDatabaseConnectionFailed
	})
}

// HasTable ...
func (m Migrator) HasTable(value interface{}) bool {
	var hasTable bool
	var err error

	err = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if dialector, ok := m.DB.Dialector.(Dialector); ok {
			hasTable, err = dialector.CollectionExists(stmt.Context, stmt.Table)
			return err
		}
		return ErrDatabaseConnectionFailed
	})
	if err != nil {
		panic(err)
	}

	return hasTable
}
