package arango

import "time"

// Config for the database.
type Config struct {
	URI                  string        // URI where to find the database server (including protocol and port).
	User                 string        // Database user name
	Password             string        // Database user password
	Database             string        // Database name to use.
	Timeout              time.Duration // Maximum duration to wait until the initial connection with the database is established.
	MaxConnectionRetries uint64        // Maximum number of connection retries.
}
