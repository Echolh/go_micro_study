package test

import (
	"fmt"
	"limiter"
	"sync"
	"testing"
	"time"
)

// 测试limiter
func TestLimiterNotEqualRates(t *testing.T) {
	// 生成速率
	limit := limiter.Every(time.Duration(1 * time.Second))

	// 生成令牌桶
	tokenBucket := limiter.New(limit, 2)

	fmt.Println("开始发起请求")
	// 模拟请求
	var succeedCount, failedCount int
	requestDuration := time.Duration(500 * time.Millisecond)
	for i := 0; i < 10; i++ {
		time.Sleep(requestDuration)
		token := tokenBucket.Allow()
		if token {
			succeedCount++
			// fmt.Printf("%d获取到令牌，执行成功！ 剩余令牌数：【%f】\n", i, math.Floor(tokenBucket.Tokens()))
			fmt.Printf("%d获取到令牌，执行成功！\n", i)

		} else {
			failedCount++
			fmt.Printf("%d未获取到令牌，执行失败.....\n", i)
		}
	}
	fmt.Printf("令牌生成速率： %v/秒\n", limit)
	fmt.Printf("请求执行速率：%v秒/次\n", requestDuration.Seconds())
	fmt.Printf("执行成功数: %v\n", succeedCount)
	fmt.Printf("执行失败数: %v\n", failedCount)
}

// 测试limiter - 令牌生成速率 = 请求速率
func TestLimiterEqualRates(t *testing.T) {
	// 生成速率
	limit := limiter.Every(time.Duration(1 * time.Second))

	// 生成令牌桶
	tokenBucket := limiter.New(limit, 2)

	fmt.Println("开始发起请求")
	// 模拟请求
	var succeedCount, failedCount int
	requestDuration := time.Duration(1 * time.Second)
	for i := 0; i < 10; i++ {
		time.Sleep(requestDuration)
		token := tokenBucket.Allow()
		if token {
			succeedCount++
			// fmt.Printf("%d获取到令牌，执行成功！ 剩余令牌数：【%f】\n", i, math.Floor(tokenBucket.Tokens()))
			fmt.Printf("%d获取到令牌，执行成功！\n", i)

		} else {
			failedCount++
			fmt.Printf("%d未获取到令牌，执行失败.....\n", i)
		}
	}
	fmt.Printf("令牌生成速率： %v个/秒\n", limit)
	fmt.Printf("请求执行速率：%v秒/次\n", requestDuration.Seconds())
	fmt.Printf("执行成功数: %v\n", succeedCount)
	fmt.Printf("执行失败数: %v\n", failedCount)
}


type testLock struct {
	mu sync.Mutex
}

// 测试死锁
func TestDeadlock(t *testing.T) {
	tl := &testLock{}
	tl.mu.Lock()
	defer tl.mu.Unlock()
	fmt.Println("TestLock")
	b(tl)

}

func b(tl *testLock) {
	// ! 重复加锁，这里会导致死锁
	tl.mu.Lock()
	defer tl.mu.Unlock()
	fmt.Println("b")
}
