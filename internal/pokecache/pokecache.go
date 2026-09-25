package pokecache

import (
	"sync"
	"time"

	"honnef.co/go/tools/lintcmd/cache"
)

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

type Cache struct {
	mu sync.Mutex
	cacheMap map[string]cacheEntry
}

func NewCache(interval time.Duration)  *Cache {
	newCacheVar :=  &Cache{
		cacheMap: make(map[string]cacheEntry),
	}
	newCacheVar.cacheMap.
	go newCacheVar.reset(interval)
	return newCacheP
}

func (c Cache) reset (interval time.Duration) {
	ticker := time.NewTicker(interval*time.Second)
	for range <-ticker.C {
		c.mu.Lock()
		for k,v := range c.cacheMap {
			if time.Since(v.createdAt) > interval {

			}
		}
		c.mu.Unlock()
	}
}

func (c Cache) Add (key string, val []byte) {

}

func (c Cache) Get (key string, val {}bytr) {

}