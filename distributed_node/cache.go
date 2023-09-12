package distributed_node

import (
	"GemmaCache/distributed_node/lru"
	"sync"
)

// Cache 缓存封装，方便后续扩展，增加其他类型的底层缓存策略：lfu、fifo等
type Cache struct {
	lru      *lru.Cache
	maxBytes int64
	mutex    sync.Mutex
}

// Add 添加 key - value
func (c *Cache) Add(key string, value *ByteView) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.lru == nil {
		c.lru = lru.NewCache(c.maxBytes, nil)
	}

	return c.lru.Add(key, value)
}

// Get 通过key获取缓存的value
func (c *Cache) Get(key string) (*ByteView, bool) {
	if key == "" {
		return nil, false
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	res, ok := c.lru.Get(key)
	if !ok {
		return nil, false
	}

	return res.(*ByteView), true
}
