package dbx

import (
	"time"
)

// Obtain the current timestamp truncated to milliseconds, which
// is the most precision most databases (including Postgres) support.
func Now() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}
