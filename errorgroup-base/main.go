package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

// 模拟获取用户信息，这个任务会成功
func fetchUserInfo(ctx context.Context) error {
	time.Sleep(1500 * time.Millisecond)

	return fmt.Errorf("获取用户信息失败")
}

// 模拟获取产品信息，这个任务会失败
func fetchProductInfo(ctx context.Context) error {
	time.Sleep(2000 * time.Millisecond)
	log.Println("获取产品信息成功")
	return nil
}

// 模拟获取库存信息，这个任务本来需要很长时间
func fetchInventoryInfo(ctx context.Context) error {
	// 这个任务需要500毫秒，但它会在200毫秒时被errgroup取消
	select {
	case <-time.After(5000 * time.Millisecond):
		log.Println("获取库存信息成功")
		return nil
	case <-ctx.Done(): // 检查上下文是否被取消
		log.Println("获取库存信息的任务被取消了，因为其他任务出错了")
		return ctx.Err()
	}
}

// 示例例子: 模拟并发获取多个数据源
// 想象一下，你的服务需要同时从三个不同的地方获取数据来完成一个请求。这三个地方分别是“用户信息”、“产品信息”和“库存信息”。其中任何一个获取失败，整个操作都算失败。
func main() {
	// 1. 创建一个 errgroup，并关联一个基础的 context
	g, ctx := errgroup.WithContext(context.Background())

	// 2. 使用 g.Go() 并发启动三个任务
	g.Go(func() error {
		return fetchUserInfo(ctx)
	})

	g.Go(func() error {
		return fetchProductInfo(ctx)
	})

	g.Go(func() error {
		return fetchInventoryInfo(ctx)
	})

	// 3. 使用 g.Wait() 等待结果
	// 它会阻塞，直到所有 g.Go() 启动的协程都返回，
	// 或者第一个错误发生。
	if err := g.Wait(); err != nil {
		log.Printf("\n--- 任务执行出错 ---\n")
		log.Printf("错误信息: %v\n", err)
	} else {
		log.Println("\n--- 所有任务都成功执行 ---")
	}
}
