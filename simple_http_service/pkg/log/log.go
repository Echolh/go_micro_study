package log

import (
	"sync"

	"go.uber.org/zap"
	// 日志轮转库
)

// zap + 日志轮转

var Logger *zap.Logger
var once sync.Once

// 初始化日志配置
func Init() {

	Logger, _ = zap.NewProduction()

	// 1. 日志轮转
	// writer:=&lumberjack.Logger{

	// }
}
