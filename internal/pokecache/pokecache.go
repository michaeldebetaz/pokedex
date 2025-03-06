package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mux      *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
		mux:      &sync.Mutex{},
	}

	go cache.reapLoop()

	return cache
}

func (c Cache) Add(key string, val []byte) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.entries[key]; ok {
		return []byte{}, false
	}

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	return val, true
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	entry, ok := c.entries[key]
	c.mux.Unlock()

	if !ok {
		return []byte{}, false
	}

	return entry.val, true
}

func (c Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		for key, entry := range c.entries {
			expiresAt := entry.createdAt.UnixMilli() + c.interval.Milliseconds()
			if expiresAt < time.Now().UnixMilli() {
				c.mux.Lock()
				delete(c.entries, key)
				c.mux.Unlock()
			}
		}
	}
}
