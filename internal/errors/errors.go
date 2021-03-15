package errors

import (
	"errors"
	"fmt"
)

var (
	ErrDatabaseConnectionFailed = errors.New("ArangoDB database connection not initialized")

	ErrMethodNotImplemented = func(entity interface{}) error {
		return fmt.Errorf("method scan not implemented for entity %v.", entity)
	}

	ErrOpeningDatabaseConnectionFailed          = errors.New("ArangoDB opening database connection failed")
	ErrOpeningDatabaseConnectionFailedWithRetry = func(retry string) error {
		return fmt.Errorf("%s. %s", ErrOpeningDatabaseConnectionFailed.Error(), retry)
	}
)
