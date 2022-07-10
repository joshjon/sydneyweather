package api

import (
	"time"
)

type valueCache[T any] struct {
	value    *T
	inserted time.Time
}

func newValueCache[T any](expiry time.Duration) *valueCache[T] {
	cache := &valueCache[T]{
		value: nil,
	}

	go func() {
		for now := range time.Tick(time.Millisecond) {
			if now.Nanosecond()-cache.inserted.Nanosecond() >= int(expiry) {
				cache.value = nil
			}
		}
	}()

	return cache
}

func (c *valueCache[T]) put(value *T) {
	c.value = value
	c.inserted = time.Now()
}

func (c *valueCache[T]) get() (*T, bool) {
	if c.value == nil {
		return nil, false
	}
	return c.value, true
}
