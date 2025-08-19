package handler

import (
	"encoding/json"
	"net/http"
	model "simple_http_svc/internal/model"
	"simple_http_svc/pkg/log"
)

var UserHandler userHandler

type userHandler struct {
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user := model.UserInfo{
		ID:   1,
		Name: "jack",
		Age:  20,
	}
	log.Logger.Error("测试错误日志")
	// panic("爆炸")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&user)
}
