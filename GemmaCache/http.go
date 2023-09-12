package GemmaCache

import (
	"GemmaCache/distributed_node/consistent_hash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/gemmacache/"
)

type HTTPPool struct {
	// 记录自身地址
	self string
	// 监听统一前缀
	basePath string

	mutex    sync.Mutex
	hashRing *consistent_hash.Map
	getters  map[string]PeerGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
		getters:  make(map[string]PeerGetter),
	}
}

// ServeHTTP
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, h.basePath) {
		http.Error(w, fmt.Sprintf("your path error.shoule has prefix %s", h.basePath), 500)
		return
	}
	h.Log("%s %s", r.Method, r.URL.Path)
	// 地址: /gemmacache/group name/key name
	parts := strings.SplitN(path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, fmt.Sprintf("your path error.shoule be %s[group name]/[key name]", h.basePath), 500)
		return
	}

	groupName := parts[0]
	keyName := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	res, err := group.Get(keyName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(res.value)
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// PickPeer 根据key在哈希环中查找缓存位置节点
func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	addr := h.hashRing.Get(key)
	h.Log("Pick peer %s", addr)
	if addr == "" || addr == h.self {
		return nil, false
	}
	return h.getters[addr], true
}

// Set 设置节点
func (h *HTTPPool) Set(addrs ...string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.hashRing = consistent_hash.NewMap(50, nil)
	h.hashRing.Add(addrs...)

	for _, addr := range addrs {
		h.getters[addr] = &HTTPGetter{basePath: addr + h.basePath}
	}
}

type HTTPGetter struct {
	basePath string
}

func (p *HTTPGetter) Get(group, key string) ([]byte, error) {
	url := fmt.Sprintf(
		"%v%v/%v",
		p.basePath,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	return ioutil.ReadAll(res.Body)
}
