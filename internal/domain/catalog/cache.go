package catalog

import (
	"sync"
	"time"
)

type cachedItem struct {
	value      any
	expiration time.Time
}

type cache struct {
	mu    sync.RWMutex
	items map[string]cachedItem
	ttl   time.Duration
}

func newCache(ttl time.Duration) *cache {
	return &cache{
		items: make(map[string]cachedItem),
		ttl:   ttl,
	}
}

func (c *cache) get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *cache) set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cachedItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	}
}

func (c *cache) invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
