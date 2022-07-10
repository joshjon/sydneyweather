package api

import (
	"time"
)

// valueCache is a cache that stores a single value in memory for the specified
// duration. A map cache is more suitable in most circumstances, however, this
// service only supports a single city which means a map is unnecessary.
// The cache value is never cleaned up in order to support retrieval of a stale
// value.
type valueCache[T any] struct {
	value      *T
	duration   time.Duration
	expiration int64
}

// newValueCache creates a new value cache.
func newValueCache[T any](duration time.Duration) *valueCache[T] {
	cache := &valueCache[T]{
		value:      nil,
		duration:   duration,
		expiration: 0,
	}

	return cache
}

// put inserts a new value into the cache. Any existing value will be overridden.
func (c *valueCache[T]) put(value *T) {
	c.value = value
	c.expiration = time.Now().Add(c.duration).UnixNano()
}

// get returns the cache value.
func (c *valueCache[T]) get() (*T, bool) {
	if c.value == nil {
		return nil, false
	}
	return c.value, true
}

// expired returns a boolean that indicates if the cache value is expired.
func (c *valueCache[T]) expired() bool {
	return c.expiration == 0 || c.expiration > 0 && time.Now().UnixNano() > c.expiration
}
