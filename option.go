package keyrot

import "time"

// Option is a configuration option for the key rotation manager.
type Option func(m *Manager)

// Duration specifies the duration for which keys will be valid.
func Duration(d time.Duration) Option {
	return func(m *Manager) {
		m.duration = d
	}
}

// Limit specifies the maximum number of keys to be retained.
// This implicitly defines the key rotation rate, which is
// duration / limit.
func Limit(limit int) Option {
	return func(m *Manager) {
		m.limit = limit
	}
}

// Bits specifies the number of bits of entropy present in the keys
// that are generated.
func Bits(bits int) Option {
	return func(m *Manager) {
		m.bits = bits
	}
}
