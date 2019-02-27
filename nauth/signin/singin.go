package signin

import "time"

// Log represent an instance of a sign-in attempt
// recorded.
type Log struct {
	UserAgent string
	At        time.Time
	IP        string
}
