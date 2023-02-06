package app_cache

import "sync"

type (
	ChecksCache struct {
		sync.RWMutex
		m map[string]string
	}

	Cache struct {
		LastChecks *ChecksCache
	}
)

func New() *Cache {
	return &Cache{
		LastChecks: &ChecksCache{
			m:     make(map[string]string),
		},
	}
}

func (c *ChecksCache) Add(key string, value string) {
	c.Lock()
	defer c.Unlock()

	c.m[key] = value
}

func (c *ChecksCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.m, key)
}

func (c *ChecksCache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()

	res, ok := c.m[key]

	return res, ok
}

func (c *ChecksCache) GetMapCopy() map[string]string {
	c.RLock()
	defer c.RUnlock()

	res := make(map[string]string)

	for key, value := range c.m {
		res[key] = value
	}

	return res
}
