package keyrot

import (
	"sync"
	"time"
)

// Manager manages time-limited rotating keys.
type Manager struct {
	duration time.Duration // Maximum validity period of a key
	limit    int           // Number of keys to maintain within the validity period
	bits     int           // Number of bits in key values (rounded to nearest byte)

	mu   sync.RWMutex
	keys []Key
}

// New returns a new manager with the given options.
func New(opts ...Option) *Manager {
	m := &Manager{
		duration: DefaultDuration,
		limit:    DefaultLimit,
		bits:     DefaultBits,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Validate returns true if the given key value is valid.
func (m *Manager) Validate(value string) bool {
	now := time.Now().Add(time.Second)
	keys := m.getKeys(now)

	for _, key := range keys {
		if key.Valid(now) && key.String() == value {
			return true
		}
	}

	return false
}

// Key returns the most recently generated key.
func (m *Manager) Key() Key {
	now := time.Now().Add(time.Second)
	keys := m.getKeys(now)
	return keys[0]
}

// getKeys will rotate the keys if necessary and return a copy of the
// current key set.
func (m *Manager) getKeys(now time.Time) (keys []Key) {
	m.mu.RLock()
	if m.shouldRotate(now) {
		m.mu.RUnlock()
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.shouldRotate(now) {
			m.rotate(now)
		}
	} else {
		defer m.mu.RUnlock()
	}

	keys = make([]Key, len(m.keys))
	copy(keys, m.keys)

	return keys
}

func (m *Manager) shouldRotate(now time.Time) bool {
	if len(m.keys) == 0 {
		return true
	}

	if !m.keys[0].Valid(now) {
		return true
	}

	if m.keys[0].Age(now) >= m.duration/time.Duration(m.limit) {
		return true
	}

	return false
}

func (m *Manager) rotate(now time.Time) {
	next := Key{
		timestamp: time.Now(), // Always use the unadjusted current time
		duration:  m.duration,
		value:     generate(m.bits),
	}
	keys := []Key{next}
	for i, key := range m.keys {
		if i+1 < m.limit && key.Valid(now) {
			keys = append(keys, key)
		}
	}
	m.keys = keys
}
