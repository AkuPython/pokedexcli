package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	m sync.Mutex
	ce map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{ce: make(map[string]cacheEntry)}

	go func() {
		for {
			time.Sleep(interval)

			cache.m.Lock()
			for k, ce := range cache.ce {
				if ce.createdAt.Add(interval).After(time.Now()) {
					delete(cache.ce, k)
				}
			}
			cache.m.Unlock()
		}
	}()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.m.Lock()
	defer c.m.Unlock()

	c.ce[key] = cacheEntry{val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	val, ok := c.ce[key]
	return val.val, ok
}
