package router

import (
	"net/http"
	"simple_http_svc/internal/config"
	user "simple_http_svc/internal/handler"
	"simple_http_svc/internal/middleware"
	"time"
)

func NewRouter(cfg *config.Config) http.Handler {

	// 创建多路复用器
	mux := http.NewServeMux()

	// 添加路由
	mux.HandleFunc("/get-user", user.UserHandler.GetUser)

	// 全局中间件
	server := middleware.Apply(mux, middleware.RequestLog(), middleware.Recover())

	// 设置超时
	handler := http.TimeoutHandler(server,
		time.Duration(cfg.Server.ReadTimeout)*time.Second,
		"request timeout")

	return handler
}
