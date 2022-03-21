package cache

import (
	"sync"
)

const (
	TypInMemory = "inmemory"
)

type inMemoryCache struct {
	m     map[string][]byte
	mutex sync.RWMutex
	Status
}

func newInMemCache() *inMemoryCache {
	return &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Status{}}
}

func (c *inMemoryCache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.m[key], nil
}

func (c *inMemoryCache) Set(key string, value []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if v, ok := c.m[key]; ok {
		c.del(key, v)
	}
	c.add(key, value)
	c.m[key] = value
	return nil
}

func (c *inMemoryCache) Del(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.m[key]; ok {
		c.del(key, v)
		delete(c.m, key)
	}
	return nil
}

func (c *inMemoryCache) GetStatus() Status {
	return c.Status
}
