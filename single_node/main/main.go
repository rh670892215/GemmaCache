package main

import (
	"GemmaCache/single_node"
	"fmt"
	"net/http"
)

/*
$ curl http://localhost:9999/gemmacache/scores/Tom
630

$ curl http://localhost:9999/gemmacache/scores/kkk
kkk not exist
*/

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	single_node.NewGroup("scores", 2<<10, single_node.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), nil)

	addr := "localhost:9999"
	peer := single_node.NewHTTPPool(addr)
	http.ListenAndServe(addr, peer)

}
