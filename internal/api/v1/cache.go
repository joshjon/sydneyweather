package api

import (
	"sync"
	"time"
)

// valueCache is a cache that stores a single value in memory for the specified
// duration. A map cache is more suitable in most circumstances, however, this
// service only supports a single city which means a map is unnecessary.
type valueCache[T any] struct {
	sync.Mutex
	value      *T
	duration   time.Duration
	expiration int64
}

// newValueCache creates a new value cache and starts a new go routine that is
// responsible for resetting the cache value once the expiry duration has elapsed.
func newValueCache[T any](duration time.Duration) *valueCache[T] {
	cache := &valueCache[T]{
		value:      nil,
		duration:   duration,
		expiration: 0,
	}

	go func() {
		for {
			if cache.expired() {
				cache.value = nil
				cache.expiration = 0
			}
		}
	}()

	return cache
}

// put inserts a new value into the cache. Any existing value will be overridden.
func (c *valueCache[T]) put(value *T) {
	c.Lock()
	defer c.Unlock()
	c.value = value
	c.expiration = time.Now().Add(c.duration).UnixNano()
}

// get returns the cache value.
func (c *valueCache[T]) get() (*T, bool) {
	c.Lock()
	defer c.Unlock()
	if c.value == nil {
		return nil, false
	}
	return c.value, true
}

// expired returns a boolean that indicates if the cache value is expired.
func (c *valueCache[T]) expired() bool {
	c.Lock()
	defer c.Unlock()
	return c.expiration > 0 && time.Now().UnixNano() > c.expiration
}
