package main

import (
	"context"
	"fmt"
	"log"
	"simple_rpc_svc/internal/config"
	"simple_rpc_svc/internal/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc客户端
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置文件失败，err%+v", err)
	}

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("无法连接到grpc服务器，err;%v", err)
	}

	defer conn.Close()
	client := proto.NewUserServiceClient(conn)

	// 控制grpc的超时时间，
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := client.GetUser(ctx, &proto.GetUserRequest{})
	if err != nil {
		fmt.Printf("获取用户失败，err%+v", err)
	}
	fmt.Printf("获取用户信息成功，user:%+v", user)

}
