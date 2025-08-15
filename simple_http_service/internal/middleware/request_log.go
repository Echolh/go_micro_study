package middleware

import (
	"log"
	"net/http"
	"simple_http_svc/internal/model"
	"time"
)

// type test func(string) string

// func NewTest() test {
// 	return func(s string) string {
// 		return s + "hello"
// 	}
// }

// 请求日志中间件
func RequestLog() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Printf(
				"started %s %s",
				r.Method,
				r.URL.Path,
			)

			// 使用ResponseWriter包装器捕获状态码
			lrw := &model.Response{ResponseWriter: w, StatusCode: http.StatusOK}

			// 把请求“传递”给下一个handler
			next.ServeHTTP(lrw, r)

			log.Printf(
				"completed %s %s %d %s",
				r.Method,
				r.URL.Path,
				lrw.StatusCode,
				time.Since(start),
			)
		})
	}
}
