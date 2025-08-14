package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"simple_rpc_svc/internal/config"
	"simple_rpc_svc/internal/server"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败,err:%+v", err)
	}

	// 创建并启动gPRC服务器
	grpcServer := server.NewGRPCServer(cfg)

	// 在goroutine中启动服务器，避免阻塞
	// grpcServer.Start() 会一直阻塞（监听端口、处理请求），
	// 用 goroutine 包裹后，主线程不会被卡住，能继续执行后面的 “监听信号” 逻辑。
	go func() {
		if err := grpcServer.Start(); err != nil {
			fmt.Printf("服务启动失败；%v", err)
		}
	}()

	// 优雅退出
	// 创建一个channel接收系统信号
	signChan := make(chan os.Signal, 1)
	// 告诉系统：当收到 Ctrl+C（SIGINT）或 kill 命令（SIGTERM）时，把信号发送到 sigChan
	// signal.Notify 是 Go 标准库的 “信号监听器”
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞等待：程序会卡在这行，直到收到终止信号
	<-signChan

	// 收到结束信号
	fmt.Print("收到终止信号，开始关闭服务器....")

	// 优雅关闭的逻辑
	// 创建一个带超时的上下文（最多等 5 秒处理剩余请求）
	// 给服务关闭设置一个 “最大等待时间”（5 秒）：如果 5 秒内所有请求都处理完了，就正常关闭；如果超时还有未处理完的请求，就强制终止（避免无限等待）。
	// defer cancel() 确保超时后释放上下文资源，防止内存泄漏
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ReadTimeout)*time.Second)
	defer cancel()
	// 调用服务的停止方法（内部会处理：停止接收新请求 + 等待现有请求完成）
	grpcServer.Stop()

	<-ctx.Done()
	log.Print("GRPC服务已完全关闭。")
}
