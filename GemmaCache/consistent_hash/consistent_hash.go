package consistent_hash

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type Hash func([]byte) uint32

type Map struct {
	hash      Hash
	hashMap   map[int]string
	replicas  int
	hashCodes []int
}

func NewMap(replicas int, hash Hash) *Map {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}

	return &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (m *Map) Add(addrs ...string) {
	for _, addr := range addrs {
		for i := 0; i < m.replicas; i++ {
			hashCode := int(m.hash([]byte(fmt.Sprintf("%d%s", i, addr))))
			m.hashMap[hashCode] = addr
			m.hashCodes = append(m.hashCodes, hashCode)
		}
	}
	sort.Ints(m.hashCodes)
}

func (m *Map) Get(key string) string {
	if key == "" {
		return ""
	}
	hashCode := int(m.hash([]byte(key)))
	// 找到第一个大于等于hashCode的addr
	index := sort.Search(len(m.hashCodes), func(i int) bool {
		return m.hashCodes[i] >= hashCode
	})

	// 取余计算，形成环
	return m.hashMap[m.hashCodes[index%len(m.hashCodes)]]
}
