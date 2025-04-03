package cache

import (
	"fmt"
	"sync"
	"time"
)

type item[V any] struct {
	value      V
	expiration time.Time
}

func (i *item[V]) isExpired() bool {
	return time.Now().After(i.expiration)
}

type cache[K fmt.Stringer, V any] struct {
	mu       sync.Mutex
	items    map[string]item[V]
	duration time.Duration
}

func newCache[K fmt.Stringer, V any](duration, gcInterval time.Duration) *cache[K, V] {
	cache := &cache[K, V]{
		items:    make(map[string]item[V]),
		duration: duration,
	}

	go cache.schedule(gcInterval)
	return cache
}

func (c *cache[K, V]) schedule(interval time.Duration) {
	for range time.NewTicker(interval).C {
		c.gc()
	}
}

func (c *cache[K, V]) gc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.isExpired() {
			delete(c.items, key)
		}
	}
}

func (c *cache[K, V]) put(key K, value V) {
	c.items[key.String()] = item[V]{
		value:      value,
		expiration: time.Now().Add(c.duration),
	}
}

func (c *cache[K, V]) get(key K) (value V, exists bool) {
	item, ok := c.items[key.String()]
	if !ok || item.isExpired() {
		return
	}

	return item.value, true
}
