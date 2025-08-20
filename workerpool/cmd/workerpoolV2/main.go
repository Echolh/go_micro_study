package workerpoolv2

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// 并发任务池

// 默认配置
const (
	defaultMaxWorkerCount = 5
	defaultQueueSize      = 100
)

// 任务
type Task interface {
	Run(ctx context.Context) error
}

// 函数任务适配器
type TaskFunc func(ctx context.Context) error

// 实现Run函数
func (t TaskFunc) Run(ctx context.Context) error {
	return t(ctx)
}

// 配置函数
type Option func(p *WorkerPool) error

// worker_pool任务池
type WorkerPool struct {
	ctx            context.Context    // context
	cancel         context.CancelFunc // 通知所有任务和工作协程终止运行，确保资源被正确释放
	queue          chan Task          // 任务队列
	queueSize      int                // 队列的容量
	maxWorkerCount int                // 最大woker数量

	wg            sync.WaitGroup
	closed        uint32                //  0/1:是否已关闭任务队列，关闭后禁止在添加任务
	firstErr      atomic.Pointer[error] // 第一个错误
	AddedCount    uint32                // 已添加任务数
	ExecutedCount uint32                // 已执行任务数（无论成功失败）
}

// New
func New(options ...Option) (*WorkerPool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &WorkerPool{
		ctx:            ctx,
		cancel:         cancel,
		queue:          make(chan Task, defaultQueueSize),
		queueSize:      defaultQueueSize,
		maxWorkerCount: defaultMaxWorkerCount,
	}

	// 配置
	for _, v := range options {
		if err := v(pool); err != nil {
			return nil, err
		}
	}

	// 启动pool
	pool.startPool()

	return pool, nil
}

// 设置queueSize
func (p *WorkerPool) SetQueueSize(size int) Option {
	return func(p *WorkerPool) error {
		if size < 1 {
			return errors.New("size 不能小于1")
		}

		p.queueSize = size
		p.queue = make(chan Task, size)
		return nil
	}
}

// 设置最大woker数量
func (p *WorkerPool) SetMaxWorkerCount(count int) Option {
	return func(p *WorkerPool) error {
		if count < 1 {
			return errors.New("MaxWorkerCount 不能小于1")
		}

		p.maxWorkerCount = count
		return nil
	}
}

// AddTaskFunc 添加任务
func (p *WorkerPool) AddTaskFunc(f func(ctx context.Context) error) error {
	return p.addTask(TaskFunc(f))
}

// startPool 开始执行
func (p *WorkerPool) startPool() {
	p.wg.Add(p.maxWorkerCount)
	for i := 0; i < p.maxWorkerCount; i++ {
		go func(workerID int) {
			defer p.wg.Done()
			// 执行task
			p.workerLoop()
		}(i)
	}
}

// Wait 等待所有任务完成，并返回第一个错误
func (p *WorkerPool) WaitAndClose() error {
	p.Shutdown()
	p.wg.Wait()
	return p.GetFirstError()
}

// workerLoop 工作协程循环
func (p *WorkerPool) workerLoop() {
	for {
		select {
		case task, ok := <-p.queue:
			if !ok {
				// 队列已关闭且所有任务已取出，退出
				return
			}
			// 执行task
			p.executeTask(task)
		case <-p.ctx.Done():
			// 上下文已取消，退出
			return
		}

	}
}

// executeTask 执行任务
func (p *WorkerPool) executeTask(t Task) {
	atomic.AddUint32(&p.ExecutedCount, 1)
	// recover
	defer func() {
		if err := recover(); err != nil {
			e := errors.Errorf("任务panic，err：%+v", err)
			p.setFirstError(e)
		}
	}()

	// 执行任务
	if err := t.Run(p.ctx); err != nil {
		p.setFirstError(err)
	}
}

func (p *WorkerPool) addTask(t Task) error {
	// time.Sleep(10 * time.Microsecond)
	// 判断是否已出错
	if err := p.GetFirstError(); err != nil {
		return errors.Errorf("任务池已出错，%v", err)
	}

	// 判断是否已关闭
	if p.IsClosed() {
		return errors.New("任务池已关闭，不允许继续添加任务！")
	}

	// 尝试向队列中添加任务
	select {
	case p.queue <- t:
		atomic.AddUint32(&p.AddedCount, 1)
		return nil
	case <-p.ctx.Done(): // 若已取消，直接返回错误
		return errors.New("任务池已取消，无法添加任务")
	default:
		return errors.New("任务队列已满")
	}
}

// setFirstError
func (p *WorkerPool) setFirstError(err error) {
	// 仅当firstErr =nil 时才设置错误
	if p.firstErr.CompareAndSwap(nil, &err) {
		p.Shutdown()
	}
}

// GetFirstError
func (p *WorkerPool) GetFirstError() error {
	if errPtr := p.firstErr.Load(); errPtr != nil {
		return *errPtr
	}
	return nil
}

func (p *WorkerPool) Shutdown() {
	// 原子操作，关闭queue
	if atomic.CompareAndSwapUint32(&p.closed, 0, 1) {
		close(p.queue)
	}

	// 出错时，取消所有任务
	if p.GetFirstError() != nil {
		// 取消所有任务
		p.cancel()
	}
}

// queue是否已关闭
func (p *WorkerPool) IsClosed() bool {
	return atomic.LoadUint32(&p.closed) == 1
}
