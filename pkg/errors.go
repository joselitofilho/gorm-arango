package arango

import (
	"errors"
	"fmt"
)

var (
	ErrDatabaseConnectionFailed                 = errors.New("ArangoDB database connection not initialized")
	ErrOpeningDatabaseConnectionFailed          = errors.New("ArangoDB opening database connection failed")
	ErrOpeningDatabaseConnectionFailedWithRetry = func(retry string) error {
		return fmt.Errorf("%s. %s", ErrOpeningDatabaseConnectionFailed.Error(), retry)
	}
)
