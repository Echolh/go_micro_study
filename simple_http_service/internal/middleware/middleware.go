package middleware

import "net/http"

// 定义中间件函数类型：接收一个处理器，返回一个新的处理器
type Middleware func(http.Handler) http.Handler

// 添加中间件
func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
