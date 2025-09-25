package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"time"
)

// 大部分任务需要 800ms 才能完成。
//
// 有一个“慢任务”（比如第8个任务），需要 2500ms 才能完成。
//
// 有一个“会出错的任务”（比如第15个任务），它会在 300ms 后返回一个错误
func processTask(ctx context.Context, taskId int) error {
	duration := time.Duration(700+rand.Intn(200)) * time.Millisecond
	// 有一个“慢任务”（比如第8个任务），需要 2500ms 才能完成。
	if taskId == 8 {
		duration = 2500 * time.Millisecond
	}
	if taskId == 15 {
		duration = 300 * time.Millisecond
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("任务%d在执行前接到取消信号", taskId)
			return fmt.Errorf("任务%d在执行前接到取消信号", taskId)
		case <-time.After(duration):
			if taskId == 15 {
				return fmt.Errorf("任务 %d 执行失败", taskId)
			}
			log.Printf("任务 %d 执行完成", taskId)
			return nil
		}
	}
}

func main() {
	// 1.定义最大并发数
	const maxConcurrency = 4
	const taskCount = 20

	// 创建一个基础的 context
	parentCtx := context.Background()

	// 1. 设定一个 2 秒的超时期限
	// 我们从 parentCtx 创建一个带有 2 秒超时的子 context
	// Go 语言的惯例是，把 cancel 函数用 defer 调用，确保无论函数如何退出，资源都会被清理
	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	// 2. 创建一个容量为 maxConcurrency 的带缓冲 channel 作为信号量
	semaphore := make(chan struct{}, maxConcurrency)

	g, ctx := errgroup.WithContext(ctx)

	log.Printf("开始执行任务")
	for i := 0; i < taskCount; i++ {
		// 在循环的顶部捕获变量，防止闭包问题
		taskId := i

		// 3. 在启动 Goroutine 之前，先获取一个“令牌”
		// 这个操作会向 channel 发送一个空结构体。如果 channel 满了（意味着并发数已达上限），
		// 这里会阻塞，直到有其他 Goroutine 完成任务并释放了令牌。
		semaphore <- struct{}{}
		g.Go(func() error {
			// 4. 使用 defer 来确保任务完成后，无论成功还是失败，都释放令牌。
			// 这个操作会从 channel 接收一个值，从而为其他等待的 Goroutine 腾出空间。
			defer func() { <-semaphore }()

			return processTask(ctx, taskId)
		})
	}

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
	log.Println("全部执行完成")
}
