package middleware

import (
	"net/http"
	"simple_http_svc/internal/config"
)

// 认证中间件
func Auth(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 逻辑
			token := r.Header.Get("Authorization")
			if token != "Bearer"+cfg.Auth.SecretToken {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// 把请求“传递”给下一个handler
			next.ServeHTTP(w, r)
		})
	}
}
