package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu       sync.Mutex
	cacheMap map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	newCacheP := &Cache{
		cacheMap: make(map[string]cacheEntry),
	}
	go newCacheP.readLoop(interval)
	return newCacheP
}

func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.cacheMap {
			if time.Since(v.createdAt) > interval {
				delete(c.cacheMap, k)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(k string, v []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[k] = cacheEntry{
		createdAt: time.Now(),
		val:       v,
	}
}

func (c *Cache) Get(k string) (v []byte, isExist bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cacheMap[k]
	return entry.val, ok
}
