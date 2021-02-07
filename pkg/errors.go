package arango

import (
	"errors"
	"fmt"
)

var (
	ErrOpeningDatabaseConnectionFailed          = errors.New("ArangoDB opening database connection failed")
	ErrOpeningDatabaseConnectionFailedWithRetry = func(retry string) error {
		return fmt.Errorf("%s. %s", ErrOpeningDatabaseConnectionFailed.Error(), retry)
	}
)
