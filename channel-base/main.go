package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"time"
)

// 模拟爬取一个 URL，需要随机耗时
func crawlURL(ctx context.Context, url string) error {
	log.Printf("开始爬取 %s", url)

	// 模拟随机的网络延迟
	crawlTime := time.Duration(500+rand.Intn(500)) * time.Millisecond

	select {
	case <-time.After(crawlTime):
		log.Printf("完成爬取 %s, 耗时 %v", url, crawlTime)
		return nil
	case <-ctx.Done():
		log.Printf("爬取 %s 的任务被取消", url)
		return ctx.Err()
	}
}

func main() {
	// 我们要爬取 10 个 URL
	urlsToCrawl := []string{
		"http://example.com/page-1",
		"http://example.com/page-2",
		"http://example.com/page-3",
		"http://example.com/page-4",
		"http://example.com/page-5",
		"http://example.com/page-6",
		"http://example.com/page-7",
		"http://example.com/page-8",
		"http://example.com/page-9",
		"http://example.com/page-10",
	}

	g, ctx := errgroup.WithContext(context.Background())

	// 1. 定义最大并发数
	const maxConcurrency = 5
	// 2. 创建一个容量为 maxConcurrency 的带缓冲 channel 作为信号量
	semaphore := make(chan struct{}, maxConcurrency)

	for _, url := range urlsToCrawl {
		// 在循环的顶部捕获变量，防止闭包问题
		url := url

		// 3. 在启动 Goroutine 之前，先获取一个“令牌”
		// 这个操作会向 channel 发送一个空结构体。如果 channel 满了（意味着并发数已达上限），
		// 这里会阻塞，直到有其他 Goroutine 完成任务并释放了令牌。
		semaphore <- struct{}{}
		g.Go(func() error {
			// 4. 使用 defer 来确保任务完成后，无论成功还是失败，都释放令牌。
			// 这个操作会从 channel 接收一个值，从而为其他等待的 Goroutine 腾出空间。
			defer func() { <-semaphore }()

			return crawlURL(ctx, url)
		})
	}

	// 等待所有爬取任务完成
	if err := g.Wait(); err != nil {
		fmt.Printf("\n--- 爬虫执行出错: %v ---\n", err)
	} else {
		fmt.Println("\n--- 所有 URL 都已成功爬取 ---")
	}
}
