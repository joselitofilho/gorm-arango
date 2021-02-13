package arango

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
	gormCallbacks "gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"

	"github.com/joselitofilho/gorm/driver/arango/internal/callbacks"
	"github.com/joselitofilho/gorm/driver/arango/internal/conn"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/cenkalti/backoff/v4"
	"github.com/sirupsen/logrus"
)

// DriverName is the default driver name for ArangoDB.
const DriverName = "gorm-arango"

// Dialector GORM ArangoDB dialector
type Dialector struct {
	DriverName string
	Config     *Config
	Conn       gorm.ConnPool

	Connection driver.Connection
	Client     driver.Client
	Database   driver.Database
}

// Open creates the dialect based on ArangoDB configuration.
func Open(config *Config) gorm.Dialector {
	return &Dialector{Config: config}
}

// Name returns the ArangoDB dialector name.
func (dialector Dialector) Name() string {
	return "arango"
}

// DatabaseExists checks if a database exists.
func (dialector Dialector) DatabaseExists(ctx context.Context, databaseName string) (bool, error) {
	var err error
	exists := false

	databases, err := dialector.Client.Databases(ctx)
	if err != nil {
		return false, err
	}

	for _, d := range databases {
		if d.Name() == databaseName {
			exists = true
		}
	}

	return exists, nil
}

// CreateDatabaseIfNeeded creates a database if it doesn't exist.
func (dialector Dialector) CreateDatabaseIfNeeded(ctx context.Context, databaseName string) (driver.Database, error) {
	var err error
	var database driver.Database

	exists, err := dialector.DatabaseExists(ctx, databaseName)
	if exists {
		database, err = dialector.Client.Database(ctx, databaseName)
	} else {
		database, err = dialector.Client.CreateDatabase(ctx, databaseName, nil)
	}
	if err != nil {
		return nil, err
	}

	return database, nil
}

// Initialize database based on dialector.Config.
func (dialector Dialector) Initialize(db *gorm.DB) error {
	ctx, cancel := dialector.setupContext()
	defer cancel()

	if dialector.DriverName == "" {
		dialector.DriverName = DriverName
	}

	logEntry := logrus.WithFields(logrus.Fields{
		"uri":      dialector.Config.URI,
		"user":     dialector.Config.User,
		"database": dialector.Config.Database,
	})
	logEntry.Debug("Connecting to ArangoDB server...")

	connection, err := http.NewConnection(http.ConnectionConfig{Endpoints: []string{dialector.Config.URI}})
	if err != nil {
		logEntry.WithError(err).Error("ArangoDB connection creation failed")
		return err
	}
	dialector.Connection = connection

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     connection,
		Authentication: driver.BasicAuthentication(dialector.Config.User, dialector.Config.Password),
	})
	if err != nil {
		logEntry.WithError(err).Error("ArangoDB client creation failed")
		return err
	}
	dialector.Client = client

	expBackoff := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), dialector.Config.MaxConnectionRetries)
	var database driver.Database
	operation := func() error {
		var err error
		database, err = dialector.CreateDatabaseIfNeeded(ctx, dialector.Config.Database)
		if err != nil {
			nextBackOff := expBackoff.NextBackOff()
			logEntry.WithError(err).Errorf("ArangoDB opening database connection failed. Retrying in %v...", nextBackOff)
			return ErrOpeningDatabaseConnectionFailedWithRetry(fmt.Sprintf("Retrying in %v...", nextBackOff))
		}
		return err
	}
	err = backoff.Retry(operation, expBackoff)
	if err != nil {
		logEntry.WithError(err).Error("ArangoDB opening database connection failed")
		return ErrOpeningDatabaseConnectionFailed
	}
	dialector.Database = database

	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		db.ConnPool = reflect.ValueOf(&conn.ConnPool{Connection: connection, Database: database}).Interface().(gorm.ConnPool)
	}

	callbacks.RegisterDefaultCallbacks(db, &gormCallbacks.Config{LastInsertIDReversed: true})
	// gormCallbacks.RegisterDefaultCallbacks(db, &gormCallbacks.Config{LastInsertIDReversed: true})

	db.Dialector = dialector

	return nil
}

// Migrator ...
func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   dialector,
		CreateIndexAfterCreateTable: true,
	}}}
}

// DataTypeOf ...
func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	// TODO: Implement
	return string(field.DataType)
}

// DefaultValueOf ...
func (dialector Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	// TODO: Implement
	return clause.Expr{SQL: ""}
}

// BindVarTo ...
func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	// TODO: Implement
	writer.WriteString("@%s")
}

// QuoteTo ...
func (dialector Dialector) QuoteTo(writer clause.Writer, str string) {
	// TODO: Implement
}

// Explain ...
func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	// TODO: Implement
	return ""
}

func (dialector Dialector) setupContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(dialector.Config.Timeout*time.Second))
}
