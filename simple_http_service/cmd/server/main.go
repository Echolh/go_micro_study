package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // 导入 pprof 包
	"simple_http_svc/internal/config"
	"simple_http_svc/internal/router"
	"time"
)

func main() {

	// 读取配置文件
	config, err := config.Load()
	if err != nil {
		panic(err)
	}
	// 注册路由
	h := router.NewRouter(config)

	// 创建http服务
	server := http.Server{
		Handler:      h,
		Addr:         ":" + config.Server.Port,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

	// 单独启用一个goroutine 运行pprof服务
	go func() {
		// pprof 服务默认监听 6060 端口（可自定义）
		// if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
		if err := http.ListenAndServe(":"+config.Server.PProfPort, nil); err != nil {
			log.Fatalf("pprof 服务启动失败, err:%v", err)
		}
	}()

	// 启动服务
	log.Printf("服务运行在端口： %s ，运行模式 ：%s ", config.Server.Port, config.Env)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP 服务启动失败， err: %+v", err)
	}

	// 启动服务器（非阻塞）
	// go func() {
	// 	log.Printf("Server is running on port %s in %s mode", cfg.Server.Port, cfg.Env)
	// 	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("Listen: %s\n", err)
	// 	}
	// }()

	// 优雅关闭
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutting down server...")

	// 5秒超时关闭
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server forced to shutdown:", err)
	// }

	log.Println("服务退出！")

}
