package limiter

import (
	"math"
	"sync"
	"time"
)

// infinite rate limit; 不限制速率
const Inf = Limit(math.MaxFloat64)

// 按固定速率生成令牌
type Limit float64

// 每秒生成多少个token
// interval：时间间隔
func Every(interval time.Duration) Limit {
	if interval <= 0 {
		return Inf
	}
	return 1 / Limit(interval.Seconds())
}

// 令牌桶
type TokenBucket struct {
	mu       sync.Mutex // 所有修改 tokens、last、的操作均在 mu锁保护下进行
	capacity int        // 桶的容量
	limit    Limit      // 令牌生成速率： n个/s, 可能会有小数，如2秒生成一个：即0.5个/s，
	tokens   float64    // 当前桶中令牌的数量
	last     time.Time  // 上次生成令牌的时间
}

// 生成令牌桶
func New(limit Limit, capacity int) *TokenBucket {
	b := &TokenBucket{
		capacity: capacity,
		limit:    limit,
		tokens:   float64(capacity),
	}
	return b
}

// 是否允许执行事件
func (lim *TokenBucket) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// 是否允许在时间t执行n个事件
func (lim *TokenBucket) AllowN(t time.Time, n int) bool {
	return lim.reserveN(t, n)
}

// 核心代码
// 预定N个令牌
func (lim *TokenBucket) reserveN(t time.Time, n int) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	// 不限速直接返回
	if lim.limit == Inf {
		lim.last = t
		lim.tokens = float64(n)
		return true
	}

	// 预生成令牌
	tokens := lim.advance(t)
	tokens = tokens - float64(n)
	if tokens < 0 {
		return false
	}

	// 更新令牌数
	lim.tokens = tokens
	lim.last = t

	return true
}

// 生成的token数量
func (lim *TokenBucket) advance(t time.Time) (newTokens float64) {
	// 第一次启动时，填满桶
	if lim.last.IsZero() {
		lim.last = t
		lim.tokens = float64(lim.capacity)
		newTokens = lim.tokens
		return
	}
	lastTime := lim.last
	if t.Before(lastTime) {
		t = lastTime
	}
	seconds := t.Sub(lastTime).Seconds()
	generated := float64(lim.limit) * seconds
	newTokens = lim.tokens + generated
	if newTokens > float64(lim.capacity) {
		newTokens = float64(lim.capacity)
	}
	return
}

func (lim *TokenBucket) Tokens() float64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.tokens
}
