package main

import (
	"context"
	"log"
	"sync"
	"time"
)

// processTask 只负责处理单个任务的逻辑
func processTask(taskId int) {
	log.Printf("任务 %d: 开始执行", taskId)

	// ??? 1: 在这里为这个任务创建一个独立的、带2秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var duration time.Duration
	if taskId == 4 || taskId == 8 || taskId == 12 {
		duration = 5 * time.Second // 慢任务
	} else {
		duration = 1 * time.Second // 普通任务
	}

	// ??? 2: 在这里使用一个 select 语句
	select {
	case <-time.After(duration):
		// 任务正常完成了，打印成功日志
		log.Printf("%d 执行成功", taskId)
	case <-ctx.Done():
		// 超时了，打印超时失败日志
		log.Printf("%d 执行超时", taskId)
	}
}

func main() {
	const maxConcurrency = 3
	const taskCount = 15

	semaphore := make(chan struct{}, maxConcurrency)

	// ??? 3: 在这里创建一个 WaitGroup
	var wg sync.WaitGroup

	for i := 1; i <= taskCount; i++ {
		semaphore <- struct{}{}

		// ??? 4: 在这里为 WaitGroup 增加计数
		wg.Add(1)

		taskId := i
		go func() {
			// ??? 5: 在这里使用 defer 来确保 WaitGroup 计数会减少，并且信号量会被释放
			defer func() {
				<-semaphore
				wg.Done()
			}()

			processTask(taskId)
		}()
	}

	// ??? 6: 在这里等待所有任务都完成
	wg.Wait()

	log.Println("所有任务的生命周期都已结束。")
}
