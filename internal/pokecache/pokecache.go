package pokecache

import (
	"fmt"
	"time"
	"sync"
)

var debug bool = false

type Cache struct {
	// Must lock this before accessing data
	mu sync.RWMutex

	// Cached byte data keyed on a string
	data map[string]cacheEntry
}

type cacheEntry struct {
	// When the entry was created
	createdAt time.Time

	// TODO if we're really nice we'd store the "cache-control: max-age"
	// response header value here and use that where possible

	// The raw data we're cachine
	val []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{}
	cache.data = make(map[string]cacheEntry)
	
	go cache.reapLoop(interval)

	return &cache
}

// Must run this inside a go func
func (c *Cache) reapLoop(interval time.Duration) {
	// This will run forever, we don't bother with a shutdown channel yet.
	reapTicker := time.NewTicker(interval)
	for now := range reapTicker.C {
		c.mu.Lock()

		for key, entry := range c.data {
			expiry := entry.createdAt.Add(interval)
			if expiry.Before(now) {
				if debug {
					fmt.Printf("Expiring cache entry for %s\n", key)
				}
				delete(c.data, key)
			}
		}

		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if debug {
		fmt.Printf("Cache HIT for %s\n", key)
	}
	return entry.val, true
}