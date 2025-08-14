package router

import (
	"net/http"
	"simple_http_svc/internal/config"
	user "simple_http_svc/internal/handler"
	"time"
)

func New(cfg *config.Config) http.Handler {

	// 创建多路复用器
	mux := http.NewServeMux()

	// 添加路由
	mux.HandleFunc("/get-user", user.UserHandler.GetUser)

	// TODO:全局中间件

	handler := http.TimeoutHandler(mux,
		time.Duration(cfg.Server.ReadTimeout)*time.Second,
		"request timeout")

	return handler
}
