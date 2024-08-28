package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache struct holds the cache entries and manages cache cleanup
type Cache struct {
	Cache      map[string]CacheEntry
	interval   time.Duration
	mu         sync.Mutex
	stopTicker chan bool
}

// CacheEntry struct holds individual cache objects and their creation time
type CacheEntry struct {
	createdAt time.Time
	object    []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Cache:      make(map[string]CacheEntry),
		interval:   interval * time.Second,
		stopTicker: make(chan bool),
	}

	// Start the reapLoop in a separate goroutine
	go cache.reapLoop()

	return cache
}

// Add inserts a new object into the cache with the current timestamp
func (c *Cache) Add(key string, object []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Cache[key] = CacheEntry{
		createdAt: time.Now(),
		object:    object,
	}
}

// Get retrieves an object from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.Cache[key]
	if !ok {
		return nil, false
	}
	return entry.object, true
}

// reapLoop runs periodically to remove outdated entries
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	for {
		select {
		case <-ticker.C:
			c.removeOldEntries()
		case <-c.stopTicker:
			ticker.Stop()
			return
		}
	}
}

// removeOldEntries removes cache entries older than the specified interval
func (c *Cache) removeOldEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.Cache {
		if now.Sub(entry.createdAt) > c.interval {
			fmt.Println("Removing old cache entry:", key)
			delete(c.Cache, key)
		}
	}
}

// Stop stops the reapLoop when the cache is no longer needed
func (c *Cache) Stop() {
	c.stopTicker <- true
}