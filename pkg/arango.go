package arango

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"

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

	// if dialector.Conn != nil {
	// 	db.ConnPool = dialector.Conn
	// } else {
	// 	db.ConnPool, _ = sql.Open("", "")
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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

	db.Dialector = Dialector{
		DriverName: dialector.DriverName,
		Config:     dialector.Config,
		Connection: connection,
		Client:     client,
		Database:   database,
	}

	return nil
}

// CollectionExists checks if a collection exists.
func (dialector Dialector) CollectionExists(collectionName string) (bool, error) {
	ctx, cancel := dialector.setupContext()
	defer cancel()

	var err error
	exists := false

	if dialector.Database == nil {
		return false, ErrDatabaseConnectionFailed
	}

	collections, err := dialector.Database.Collections(ctx)
	if err != nil {
		return false, err
	}

	for _, c := range collections {
		if c.Name() == collectionName {
			exists = true
		}
	}

	return exists, nil
}

// CreateCollection ...
func (dialector Dialector) CreateCollection(name string) (driver.Collection, error) {
	ctx, cancel := dialector.setupContext()
	defer cancel()

	if dialector.Database == nil {
		return nil, ErrDatabaseConnectionFailed
	}

	return dialector.Database.CreateCollection(ctx, name, &driver.CreateCollectionOptions{})
}

// Migrator ...
func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	// TODO: Implement
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
	if field.AutoIncrement {
		return clause.Expr{SQL: "autoincrement"}
	}
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
