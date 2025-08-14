package server

import (
	"fmt"
	"net"
	"simple_rpc_svc/internal/config"
	"simple_rpc_svc/internal/proto"
	"simple_rpc_svc/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCServer gPRC服务器
type GRPCServer struct {
	cfg     *config.Config
	Service *service.UserService
	Server  *grpc.Server
}

// NewGRPCServer 创建gRPC服务器实例
func NewGRPCServer(cfg *config.Config) *GRPCServer {
	// 创建用户服务
	userService := service.NewUserService()

	// 创建gPRC服务器
	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()), // 开发环境使用不安全的凭据
	)

	// 注册服务
	// userService 实现了GetUser函数，所以实现了UserServiceServer接口
	proto.RegisterUserServiceServer(grpcServer, userService)

	return &GRPCServer{
		cfg:     cfg,
		Service: userService,
		Server:  grpcServer,
	}
}

// Start 启动服务
func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("监听失败，%s", addr)
	}
	fmt.Printf("gRPC服务器启动，监听地址: %s\n", addr)
	if err := s.Server.Serve(listen); err != nil {
		fmt.Printf("服务启动失败，err:%+v", err)
	}
	return nil
}

// Stop 优雅关闭服务器
func (s *GRPCServer) Stop() {
	fmt.Println("开始优雅关闭gRPC服务器...")
	s.Server.GracefulStop()
	fmt.Println("gRPC服务器已关闭")
}
