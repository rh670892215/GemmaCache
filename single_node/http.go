package single_node

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	defaultBasePath = "/gemmacache/"
)

type HTTPPool struct {
	// 记录自身地址
	self string
	// 监听统一前缀
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// ServeHTTP
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, h.basePath) {
		http.Error(w, fmt.Sprintf("your path error.shoule has prefix %s", h.basePath), 500)
		return
	}

	// 地址: /gemmacache/
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
