package cache

import (
	"sync"
	"time"
)

const (
	// No Expiration
	DefaultExpiration int64 = 0
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, v interface{}, expiration time.Duration)
	Delete(key string)
}

type item struct {
	Value      interface{}
	Expiration int64
}

type inMemoryCache struct {
	lock   sync.Mutex
	values map[string]*item
}

func New() *inMemoryCache {
	return &inMemoryCache{values: make(map[string]*item)}
}

func (c *inMemoryCache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	v, found := c.values[key]
	if !found {
		return nil, false
	}

	if v.Expiration > 0 {
		if time.Now().UnixNano() > v.Expiration {
			delete(c.values, key)
			return nil, false
		}
	}

	return v.Value, true
}

func (c *inMemoryCache) Set(key string, v interface{}, expiration time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	exp := DefaultExpiration
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	c.values[key] = &item{
		Value:      v,
		Expiration: exp,
	}
}

func (c *inMemoryCache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.values, key)
}
