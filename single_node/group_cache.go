package single_node

import (
	"GemmaCache/single_node/lru"
	"errors"
	"fmt"
	"sync"
)

var (
	groups map[string]*Group
	mutex  sync.RWMutex
)

type Group struct {
	mainCache *Cache
	name      string
	getter    Getter
}

// 初始化map
func init() {
	if groups == nil {
		groups = make(map[string]*Group)
	}
}

// NewGroup 新建group
func NewGroup(name string, maxBytes int64, getter Getter, callback func(key string, value lru.Value)) *Group {
	group := &Group{
		name:   name,
		getter: getter,
		mainCache: &Cache{
			maxBytes: maxBytes,
			lru:      lru.NewCache(maxBytes, callback),
		},
	}

	mutex.Lock()
	defer mutex.Unlock()
	groups[name] = group
	return group
}

// GetGroup 根据name在groups获取group
func GetGroup(name string) *Group {
	mutex.RLock()
	defer mutex.RUnlock()

	res, ok := groups[name]
	if !ok {
		return nil
	}

	return res
}

// Get 根据key获取缓存中的value
func (g *Group) Get(key string) (*ByteView, error) {
	res, ok := g.mainCache.Get(key)
	if ok {
		return res, nil
	}

	return g.load(key)
}

// 根据用户自定义的操作去数据库中查询指定的key
func (g *Group) load(key string) (*ByteView, error) {
	// 若缓存中查询不到，则根据用户注册的查询方法去数据库中查询数据，并更新缓存
	dbRes, err := g.getter.Get(key)
	if err != nil {
		return nil, err
	}

	byteView := &ByteView{value: dbRes}
	if ok := g.mainCache.Add(key, byteView); !ok {
		return nil, errors.New(fmt.Sprintf("group name is %s,mainCache add error,key :%s,value: %s",
			g.name, key, byteView.String()))
	}

	return byteView, nil
}
