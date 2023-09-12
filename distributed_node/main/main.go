package main

import (
	"GemmaCache/distributed_node"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	// 解析输入参数
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "GemmaCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, addr := range addrMap {
		addrs = append(addrs, addr)
	}

	group := distributed_node.NewGroup("scores", 2<<10, distributed_node.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), nil)

	apiAddr := "http://localhost:9999"
	if api {
		go func() {
			http.Handle("/api", http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					key := r.URL.Query().Get("key")
					view, err := group.Get(key)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/octet-stream")
					w.Write(view.ByteSlice())

				}))
			http.ListenAndServe(apiAddr[7:], nil)
		}()
	}
	startCacheServer(addrMap[port], addrs, group)
}

func startCacheServer(addr string, addrs []string, group *distributed_node.Group) {
	peers := distributed_node.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("gemmacache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}
