package middleware

import "net/http"

// 中间件结构体
type Middleware func(http.Handler) http.Handler

// 应用中间件

func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
