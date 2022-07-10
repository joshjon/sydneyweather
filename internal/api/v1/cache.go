package api

import (
	"sync"
	"time"
)

type valueCache[T any] struct {
	sync.Mutex
	value      *T
	duration   time.Duration
	expiration int64
}

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

func (c *valueCache[T]) put(value *T) {
	c.Lock()
	defer c.Unlock()
	c.value = value
	c.expiration = time.Now().Add(c.duration).UnixNano()
}

func (c *valueCache[T]) get() (*T, bool) {
	c.Lock()
	defer c.Unlock()
	if c.value == nil {
		return nil, false
	}
	return c.value, true
}

func (c *valueCache[T]) expired() bool {
	c.Lock()
	defer c.Unlock()
	return c.expiration > 0 && time.Now().UnixNano() > c.expiration
}
