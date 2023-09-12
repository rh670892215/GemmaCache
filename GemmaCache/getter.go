package GemmaCache
/*
接口型函数，原理参考官方包中的HandlerFunc，源码如下:

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}

使用场景:
func main() {
	http.HandleFunc("/home", home)
	_ = http.ListenAndServe("localhost:8000", nil)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}

可以看到 HandleFunc 的第二个参数是 HandlerFunc 类型，此时 HandleFunc 函数的第二个参数可以传入:
1. 实现了 ServeHTTP 方法的结构体;
2. 函数签名是 type HandlerFunc 的匿名函数;
3. 函数签名是 type HandlerFunc 的函数;

使用场景:
这种方式适用于逻辑较为复杂的场景，如果对数据库的操作需要很多信息，地址、用户名、密码，还有很多中间状态需要保持，比如超时、重连、加锁等等。
这种情况下，更适合封装为一个结构体作为参数。这样，既能够将普通的函数类型（需类型转换）作为参数，也可以将结构体作为参数，使用更为灵活。
*/
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 接口型函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

