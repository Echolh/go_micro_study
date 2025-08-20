package test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	workerpoolv2 "workerpool/cmd/workerpoolV2"
)

func TestNewPool(t *testing.T) {
	pool, err := workerpoolv2.New()
	if err != nil {
		fmt.Printf("初始化任务池失败，err:%+v", err)
	}
	pool.SetMaxWorkerCount(10)
	pool.SetQueueSize(100)

	for i := 0; i < 100; i++ {
		_ = pool.AddTaskFunc(func(ctx context.Context) error {
			time.Sleep(10 * time.Microsecond)
			fmt.Println("任务" + strconv.Itoa(i))
			return nil
		})
	}

	if err := pool.WaitAndClose(); err != nil {
		fmt.Printf("WaitAndClose err:%v", err)
	}

	fmt.Printf("添加的任务数: %v\n", pool.AddedCount)
	fmt.Printf("已执行完成的任务数: %v\n", pool.ExecutedCount)

}

// 测试正常提交并执行任务
func TestSubmitAndExecute(t *testing.T) {
	pool, err := workerpoolv2.New()
	if err != nil {
		t.Fatalf("创建任务池失败: %v", err)
	}
	// defer pool.WaitAndClose()

	var wg sync.WaitGroup

	// 提交10个任务
	taskCount := 10
	wg.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		err := pool.AddTaskFunc(func(ctx context.Context) error {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond) // 模拟任务执行
			return nil
		})
		if err != nil {
			t.Errorf("提交任务 %d 失败: %v", i, err)
		}
	}

	wg.Wait()
	pool.WaitAndClose()
	if pool.GetFirstError() != nil {
		t.Errorf("预期无错误，实际得到: %v", pool.GetFirstError())
	}
	if !pool.IsClosed() {
		t.Error("Wait() 后任务池应处于关闭状态")
	}
}

// 测试第一个错误捕获
func TestFirstErrorCapture(t *testing.T) {
	pool, err := workerpoolv2.New()
	if err != nil {
		t.Fatalf("创建任务池失败: %v", err)
	}
	pool.SetMaxWorkerCount(2)
	pool.SetQueueSize(5)

	// 提交3个任务，第二个会出错
	err = pool.AddTaskFunc(func(ctx context.Context) error {
		// time.Sleep(2 * time.Millisecond)
		fmt.Println("任务1执行完成")
		//log.Logger.Info("任务1执行完成")
		return nil
	})
	if err != nil {
		t.Error("提交第一个任务失败:", err)
	}

	expectedErr := errors.New("故意出错的任务")
	err = pool.AddTaskFunc(func(ctx context.Context) error {
		// time.Sleep(2 * time.Millisecond)
		fmt.Println("任务2执行完成,随后报错")
		//log.Logger.Info("任务2执行完成")
		return expectedErr
	})
	if err != nil {
		t.Error("提交第二个任务失败:", err)
	}

	// 第三个任务应该提交失败（因为第二个任务已触发错误）
	err = pool.AddTaskFunc(func(ctx context.Context) error {
		fmt.Println("任务3执行完成")
		//log.Logger.Info("任务3执行完成")
		return nil
	})
	if err != nil {
		t.Errorf("第三个任务应提交失败,err:%v", err)
	}

	time.Sleep(2 * time.Second)
	for i := 4; i <= 13; i++ {
		err = pool.AddTaskFunc(func(ctx context.Context) error {
			fmt.Println("任务" + strconv.Itoa(i) + "执行完成")
			//log.Logger.Info("任务3执行完成")
			return nil
		})
		if err != nil {
			t.Errorf("第"+strconv.Itoa(i)+"个任务应提交失败, err:%v", err)
		}
	}

	// 等待任务池处理完成
	err1 := pool.WaitAndClose()
	if err1 != nil {
		t.Errorf("预期错误 %v，实际得到 %v", expectedErr, err1)
	}
	t.Log("添加的任务数：", pool.AddedCount)
	t.Log("执行任务数：", pool.ExecutedCount)
	// if err := pool.WaitAndClose(); err != expectedErr {
	// 	t.Errorf("预期错误 %v，实际得到 %v", expectedErr, err)
	// }
	// errN := pool.AddTaskFunc(func(ctx context.Context) error {
	// 	fmt.Println("任务N执行完成")
	// 	//log.Logger.Info("任务3执行完成")
	// 	return nil
	// })
	// if errN != nil {
	// 	t.Errorf("第N个任务提交失败，err:%v", errN)
	// }
}

// 测试任务panic捕获
func TestTaskPanicRecovery(t *testing.T) {
	pool, err := workerpoolv2.New()
	if err != nil {
		t.Fatalf("创建任务池失败: %v", err)
	}

	panicMsg := "任务内部发生panic"
	err = pool.AddTaskFunc(func(ctx context.Context) error {
		panic(panicMsg)
	})
	if err != nil {
		t.Error("提交任务失败:", err)
	}

	// 等待任务处理完成
	finalErr := pool.WaitAndClose()
	if finalErr == nil {
		t.Error("预期捕获panic错误，但未得到任何错误")
	} else if finalErr.Error() != "任务panic: "+panicMsg {
		t.Errorf("预期捕获特定panic错误，实际得到: %v", finalErr)
	}
}
