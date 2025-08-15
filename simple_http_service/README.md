# 记录编码过程


## 2025年08月14日

 使用Go原生`net/http`库 创建HTTP服务，参考[net/http简介](https://darjun.github.io/2021/07/13/in-post/godailylib/nethttp/)

![img](https://darjun.github.io/img/in-post/godailylib/nethttp1.png#center)

### 核心代码

```go
// 注册路由
	h := router.NewRouter(config)

	// 创建http服务
	server := http.Server{
		Handler:      h,
		Addr:         ":" + config.Server.Port,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
	}

 // 启动服务
	log.Printf("Server is running on port %s in %s mode", config.Server.Port, config.Env)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("http server err: %+v", err)
	}


func NewRouter(cfg *config.Config) http.Handler {

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


	
```

### 测试

![image-20250814225835298](https://data-lh.top/image-20250814225835298.png)

## 2025年08月15日

 增加中间件功能

有时候需要在请求处理代码中增加一些通用的逻辑，如统计处理耗时、记录日志、捕获宕机等等。如果在每个请求处理函数中添加这些逻辑，代码很快就会变得不可维护，添加新的处理函数也会变得非常繁琐。所以就有了中间件的需求。

中间件有点像面向切面的编程思想，但是与 Java 语言不同。在 Java 中，通用的处理逻辑（也可以称为切面）可以通过反射插入到正常逻辑的处理流程中，在 Go 语言中基本不这样做。

Go 语言中的函数是第一类值，既可以作为参数传给其他函数，也可以作为返回值从其他函数返回。我们前面介绍了处理器/函数的使用和实现。

中间件中的写法比较特殊，是函数式编程。

```go
// 定义中间件函数类型：接收一个处理器，返回一个新的处理器
type Middleware func(http.Handler) http.Handler

// 添加中间件
func Apply(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}


// http.HandlerFunc：这是 Go 语言标准库提供的一个类型，本质是 “能处理 HTTP 请求的函数”。
// 只要一个函数的参数是(w http.ResponseWriter, r *http.Request)，
// 就能被转换成http.Handler，从而作为 HTTP 服务器的处理器。

// 编写中间件
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("recovered panic: %v", err)
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			// 把请求“传递”给下一个handler
			next.ServeHTTP(w, r)
		})
	}
}

```



### 测试

![image-20250815145632538](https://data-lh.top/image-20250815145632538.png)