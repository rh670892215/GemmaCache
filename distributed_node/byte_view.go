package distributed_node

import "unsafe"

// ByteView 缓存value结构体
type ByteView struct {
	value []byte
}

// Len 返回缓存长度
func (b *ByteView) Len() int {
	return len(b.value)
}

// String string类型返回缓存value
func (b *ByteView) String() string {
	return *(*string)(unsafe.Pointer(&b.value))
}

// ByteSlice []byte形式返回缓存value
func (b *ByteView) ByteSlice() []byte {
	return cloneBytes(b.value)
}

// 返回value副本，防止在外层的修改影响缓存对象
func cloneBytes(b []byte) []byte {
	res := make([]byte, len(b))
	copy(res, b)
	return res
}
