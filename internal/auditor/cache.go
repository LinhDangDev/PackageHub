package auditor

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// Cache provides a thread-safe in-memory cache with TTL.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]cacheItem
	defaultTTL time.Duration
}

// NewCache creates an in-memory cache with the specified default TTL.
func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		items:      make(map[string]cacheItem),
		defaultTTL: defaultTTL,
	}
}

// Get retrieves an item from the cache if not expired.
func (c *Cache) Get(key string) (interface{}, bool) {
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

// Set stores an item with default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores an item with custom TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}
