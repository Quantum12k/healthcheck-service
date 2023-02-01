package app_cache

import "sync"

type (
	ChecksCache struct {
		sync.Mutex
		m map[string]string
	}

	Cache struct {
		LastChecks *ChecksCache
	}
)

func New() *Cache {
	return &Cache{
		LastChecks: &ChecksCache{
			Mutex: sync.Mutex{},
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
	c.Lock()
	defer c.Unlock()

	res, ok := c.m[key]

	return res, ok
}

func (c *ChecksCache) GetMapCopy() map[string]string {
	c.Lock()
	defer c.Unlock()

	res := make(map[string]string)

	for key, value := range c.m {
		res[key] = value
	}

	return res
}
