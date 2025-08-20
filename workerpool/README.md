      简单的并发任务池

## 简介
一个简单的并发任务池。
- 可以指定最大并发goroutine数量。
- 返回第一个error
- 生命周期管控，通过context和cancel来管控所有goroutine的生命周期，一旦一个goroutine出现error则通知所有goroutine停止


## test

### 并发执行
- 任务个数100
- 每个任务执行时间 10毫秒
串行执行的情况下， 理论执行时间应该>=1s, 通过workerpool并发执行后，只需要0.18秒
```shell
=== RUN   TestNewPool
任务0
任务2
任务3
任务1
任务7
...... # 太长了，这里省略
任务73
任务72
任务76
任务74
任务79
任务81
任务82
任务77
任务78
任务80
任务85
任务83
任务87
任务84
任务86
任务88
任务89
任务92
任务91
任务90
任务95
任务96
任务93
任务94
任务97
任务99
任务98
添加的任务数: 100
已执行完成的任务数: 100
--- PASS: TestNewPool (0.00s)
PASS
ok      workerpool/test 0.187s
```

## 心得
 在编写这个并发任务池的过程中很有收获，同时也出现了很多错误。

 ### 1. 我的生命周期，我做主

 ```golang
 // 这是我v1 版本
 func New(c context.Context, options ...Option) (*WorkerPool, error) {
    p:=&WorkerPool{
        ctx: c
    }
 }

 // 改进后的v2版本
  func New( options ...Option) (*WorkerPool, error) {

 }

 ```
这个里面犯了一个很致命的错误，就是接收外界的context传入到pool实例中，这个context可能会是致命的。按照毛大的说法，尽量不要使用requestContext来作为你的context实例，因为你不知道这个context什么时候会cancal,会导致任务池异常终止。其实还会引起很多其他问题。
1. **外部 context 可能破坏任务池的生命周期自主性** ：
任务池的核心是需要自主控制生命周期（通过 Shutdown()、Wait() 等方法主动关闭），而外部传入的 context 可能被用户在外部提前取消（例如超时、手动取消），导致：
任务池还未处理完任务就被意外终止（因为外部 context 取消会触发 ctx.Done()，工作协程退出）；
任务池的 Shutdown() 逻辑与外部 context 的取消逻辑冲突（例如外部取消后，内部 Shutdown() 可能无法正确释放资源）。
1. **职责混淆：谁来主导关闭？**
   原设计中，任务池的关闭逻辑是明确的：由内部 cancel 函数（配合 Shutdown()）主导，确保 “关闭” 操作与任务池自身状态强绑定。
   若引入外部 context，则关闭的触发源变成了两个：
   - 外部 context 的取消（如 ctx.Done()）；
   - 内部 Shutdown() 调用（触发 cancel()）。

### 2. 创建即启动，不要手动Start()

```go
//  v1 版本： Start 开始执行
func (p *WorkerPool) Start() {

}

// v2版本中，不在提供Start函数，改为在New函数中隐式调用
func (p *WorkerPool) New(options ...option)(*WorkerPool,error){
  p:=&WorkerPool{
        ctx: c
    }

   // 初始化完成后，直接启动
   p.start()
}

```
在v1版本中，我原计划是在New()之后，让调用者显示调用Start()来启动任务池。后面发现这种方式不如隐式调用，即在New()函数中直接内部调用start()函数，用户创建即启动。这种创建即启动的方式在Go生态中有很多的例子。这样做有几大好处：
1. **降低使用门槛，符合直觉**：
   对大多数用户而言，创建任务池（New()）的预期就是 “得到一个可以立即使用的任务池”，而不是 “得到一个需要手动启动的‘未激活’对象”。
例如：
当你创建一个 http.Client 时，不需要手动调用 Start() 就能立即发送请求；
当你创建一个 sync.Pool 时，也不需要启动就能直接使用。
隐式启动符合 “创建即可用” 的直觉，用户无需记住 “先 New() 再 Start()” 的额外步骤，降低了学习和使用成本。
2. **避免 “未启动” 状态的误用**：
   如果将启动逻辑拆分到 Start() 中，可能会出现用户 “忘记调用 Start()” 的情况，导致任务提交后永远不会被执行（工作协程未启动，任务队列中的任务无人处理），且这种错误难以排查（编译不会报错，运行时无明显异常，只是任务 “悄无声息地丢失”）。
在 New() 中隐式启动，可以从根源上避免这种误用 —— 任务池一旦创建，就处于 “就绪状态”，提交的任务能被立即处理（或进入队列等待），减少潜在的 bugs。

### 3. 无锁优于mutex
一开始的设计中，为了保证部分字段的线程安全，我在WorkPool中设计了一个mutex字段。后面优化为atomic包中的原子操作来取代mutex。
```golang

// v1版本，使用mutex加锁保证isClosed字段的线程安全
type WorkerPool struct {
	ctx            context.Context    // context
	cancel         context.CancelFunc // 通知所有任务和工作协程终止运行，确保资源被正确释放
	queue          chan Task          // 任务队列
	queueSize      int                // 队列的容量
	maxWorkerCount int                // 最大woker数量
	wg             sync.WaitGroup
	mutex          sync.Mutex
    isClosed bool                   // 是否关闭queue
    }

// 关闭
func (p *workerPool)Shutdown(){
     p.mutex.Lock()
    defer p.mutex.Unlock()

    p.isClosed = true
}

// 是否已关闭
func (p *workerPool) IsClosed()bool{
    p.mutex.Lock()

    defer p.mutex.Unlock()
    return p.isClosed
}

```


 ```golang
 // v2 版本使用atomic的原子操作
 type WorkerPool struct {
	ctx            context.Context    // context
	cancel         context.CancelFunc // 通知所有任务和工作协程终止运行，确保资源被正确释放
	queue          chan Task          // 任务队列
	queueSize      int                // 队列的容量
	maxWorkerCount int                // 最大woker数量
	wg             sync.WaitGroup
    isClosed unit32                   // 是否关闭queue 0/1
    }



func (p *workerPool)Shutdown(){
	// CAS原子操作
	if atomic.CompareAndSwapUint32(&p.closed, 0, 1) {
		close(p.queue)
		// 取消所有任务
	}
 
}

// 是否已关闭-atomic.LoadUint32 原子操作
func (p *workerPool) IsClosed()bool{
   return atomic.LoadUint32(&p.closed) == 1
}

 ```

 ### 4. "接口 + 适配器" 模式

- 接口：通过 Task 接口统一所有任务的func签名
- 适配器：通过 TaskFunc 适配器，用户可以直接提交普通函数作为任务，适配器会自动转换为实现了Task接口的类型
  
```golang
// 任务接口
type Task interface {
	Run(ctx context.Context) error
}

// 函数任务适配器
type TaskFunc func(ctx context.Context) error

// 实现Run函数
func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

// 添加func
func (p *workerpool)AddTaskFunc(f func(ctx context.Context)err){
      // 适配器将func转换为task的结构体
      task:=TaskFunc(f)

      // 加入队列
      p.queue<- task
}

```