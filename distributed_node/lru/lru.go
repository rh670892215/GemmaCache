package lru

import (
	"container/list"
)

type Value interface {
	Len() int
}

// Entity 缓存的实体
type Entity struct {
	key   string
	value Value
}

// Cache lru缓存
type Cache struct {
	maxBytes  int64
	usedBytes int64
	ll        *list.List
	table     map[string]*list.Element
	callBack  func(key string, value Value)
}

// NewCache 新建一个lru缓存
func NewCache(maxBytes int64, callBack func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		table:    make(map[string]*list.Element),
		ll:       list.New(),
		callBack: callBack,
	}
}

// Add 添加key-value
func (c *Cache) Add(key string, value Value) bool {
	oldVal, ok := c.table[key]
	if ok {
		kv, ok := oldVal.Value.(*Entity)
		if !ok {
			return false
		}
		c.usedBytes -= int64(kv.value.Len() - value.Len())
		c.ll.MoveToFront(oldVal)
		kv.value = value
	} else {
		entity := &Entity{key, value}
		element := c.ll.PushFront(entity)
		c.table[key] = element
		c.usedBytes += int64(len(key) + value.Len())
	}

	// 校验是否需要进行淘汰
	if c.maxBytes > 0 && c.usedBytes > c.maxBytes {
		c.delOldest()
	}
	return true
}

// Get 根据key获取指定的value
func (c *Cache) Get(key string) (Value, bool) {

	val, ok := c.table[key]
	if !ok {
		return nil, false
	}

	kv, ok := val.Value.(*Entity)
	if !ok {
		return nil, false
	}

	return kv.value, true
}

// 删除最近最少使用的元素，即队尾元素
func (c *Cache) delOldest() {

	element := c.ll.Back()
	kv, ok := element.Value.(*Entity)
	if !ok {
		return
	}
	delete(c.table, kv.key)
	c.ll.Remove(element)
	c.usedBytes -= int64(len(kv.key) + kv.value.Len())

	if c.callBack != nil {
		c.callBack(kv.key, kv.value)
	}
}

// Len lru缓存长度
func (c *Cache) Len() int {
	return len(c.table)
}
