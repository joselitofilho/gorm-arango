package arango

import (
	driver "github.com/arangodb/go-driver"
	"github.com/joselitofilho/gorm/driver/arango/internal/errors"
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
			if dialector.Database == nil {
				return errors.ErrDatabaseConnectionFailed
			}
			_, err := dialector.Database.CreateCollection(stmt.Context, stmt.Table, &driver.CreateCollectionOptions{})
			return err
		}
		return errors.ErrDatabaseConnectionFailed
	})
}

// DropTable ...
func (m Migrator) DropTable(values ...interface{}) error {
	values = m.ReorderModels(values, false)
	for i := len(values) - 1; i >= 0; i-- {
		if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
			if dialector, ok := m.DB.Dialector.(Dialector); ok {
				if hasTable := m.HasTable(stmt.Table); hasTable {
					collection, err := dialector.Database.Collection(stmt.Context, stmt.Table)
					if err != nil {
						return err
					}
					return collection.Remove(stmt.Context)
				}
				return nil
			}
			return errors.ErrDatabaseConnectionFailed
		}); err != nil {
			return err
		}
	}
	return nil
}

// HasTable ...
func (m Migrator) HasTable(value interface{}) bool {
	var hasTable bool
	var err error

	err = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if dialector, ok := m.DB.Dialector.(Dialector); ok {
			hasTable, err = dialector.Database.CollectionExists(stmt.Context, stmt.Table)
			return err
		}
		return errors.ErrDatabaseConnectionFailed
	})
	if err != nil {
		panic(err)
	}

	return hasTable
}
