package model

import "net/http"

// loggingResponseWriter 用于捕获响应状态码的ResponseWriter包装器
type Response struct {
	http.ResponseWriter
	StatusCode int
}
