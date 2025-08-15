package middleware

import (
	"log"
	"net/http"
)

// 全局recover 中间件

func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		// http.HandlerFunc：这是 Go 语言标准库提供的一个类型，本质是 “能处理 HTTP 请求的函数”。
		// 只要一个函数的参数是(w http.ResponseWriter, r *http.Request)，
		// 就能被转换成http.HandlerFunc，从而作为 HTTP 服务器的处理器。
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("recovered panic: %v", err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			// 把请求“传递”给下一个handler
			next.ServeHTTP(w, r)
		})
	}
}
