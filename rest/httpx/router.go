package httpx

import "net/http"

// Router 接口表示一个处理http请求的路由器
// Router interface represents a http router that handles http requests.
type Router interface {
	// 内嵌官方http.Handler接口
	http.Handler
	// Handle 注册一个处理http请求的路由
	// Register a route to handle http requests.
	Handle(method, path string, handler http.Handler) error
	// 设置404处理函数
	SetNotFoundHandler(handler http.Handler)
	// 设置405处理函数（请求方法不允许）
	SetNotAllowedHandler(handler http.Handler)
}
