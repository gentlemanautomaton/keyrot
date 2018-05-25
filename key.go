package keyrot

import "time"

// Key is a time-limited key with a cryptographically secure random value
// encoded as a hexadecimal string. Keys can be safely passed by value.
type Key struct {
	timestamp time.Time
	duration  time.Duration
	value     string
}

// Age returns the age of the key at the given time.
func (key Key) Age(at time.Time) time.Duration {
	if at.After(key.timestamp) {
		return at.Sub(key.timestamp)
	}
	return 0
}

// Valid returns true if the key is valid at the given time.
func (key Key) Valid(at time.Time) bool {
	return key.Age(at) <= key.duration
}

// String returns the value of the key as a hexadecimal string.
func (key Key) String() string {
	return key.value
}
