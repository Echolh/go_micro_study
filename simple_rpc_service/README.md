# 简介

```
your-project/
├── cmd/
│   ├── server/
│   │   └── main.go   # 服务端入口
│   └── client/
│       └── main.go   # 客户端入口
├── internal/
│   ├── proto/
│   │   ├── user.proto
│   │   ├── user.pb.go         # 生成的消息代码
│   │   └── user_grpc.pb.go    # 生成的gRPC代码
│   ├── service/       # 业务逻辑实现
│   └── server/        # 服务启动逻辑
├── config.yaml        # 配置文件
└── go.mod             # 依赖管理
```



## 2025年08月14日16

编写简单的rpc服务：
1. 编写.proto描述文件
2. 编译生成.pb.go文件
```shell
# 在poto目录下执行 user.proto
protoc --go_out=. --go-grpc_out=. user.proto
```
执行后，会在 proto 目录下直接生成两个文件：
plaintext
simple_rpc_svc/
└── proto/
    ├── user.proto
    ├── user.pb.go         # 消息结构体代码（在proto目录下）
    └── user_grpc.pb.go    # gRPC服务代码（在proto目录下）

3. 服务端实现约定的接口并提供服务
   

​      在internal/server 目录中 新建文件server.go, 关键函数 `NewGRPCServer()`

```go
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
```



​     客户端按照约定调用.pb.go文件中的方法请求服务



#### 启动服务

```bash
go run cmd/server/main.go
```

![image-20250814230010114](https://data-lh.top/image-20250814230010114.png)



#### 客户端调用

![image-20250814230038399](https://data-lh.top/image-20250814230038399.png)



### 为什么需要优雅关闭？

想象一个场景：用户正在通过 RPC 服务提交订单，此时你执行了 `kill` 命令关闭服务。

- 如果**没有优雅关闭**：服务会立即退出，订单请求可能只处理了一半（比如钱扣了但订单没创建），导致数据不一致。
- 如果**有优雅关闭**：服务会先拒绝新请求，等当前订单处理完再退出，确保数据正确。

总结：优雅关闭的核心是 “收到终止信号后，先处理完手头的工作，再安全退出”，这段代码通过 “信号监听 + 超时控制 + 框架的优雅关闭方法” 实现了这个目标。
