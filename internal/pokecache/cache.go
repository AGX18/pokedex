package pokecache

import (
	"sync"
	"time"
)

// This package will be responsible for all of our caching logic.

type Cache struct {
	entries map[string]cacheEntry
	mux     *sync.RWMutex // Changed to RWMutex for better read performance
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	var cache = &Cache{
		entries: make(map[string]cacheEntry),
		mux:     &sync.RWMutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // prevent ticker leak

	for range ticker.C {
		c.mux.Lock()
		for k, entry := range c.entries {
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, k)
			}
		}
		c.mux.Unlock()
	}
}
