package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// 模拟调用一个 API，它需要 5 秒才能返回
func callSlowAPI(ctx context.Context) (string, error) {
	fmt.Println("开始调用 API，预计耗时 5 秒...")

	select {
	case <-time.After(5 * time.Second):
		// 如果能走到这里，说明没有超时
		return "成功获取API数据", nil
	case <-ctx.Done():
		// 在等待的 5 秒内，如果 context 被取消了，就会走到这里
		fmt.Println("API 调用被中断!")
		return "", ctx.Err() // ctx.Err() 会返回导致取消的原因
	}
}

func main() {
	// 创建一个基础的 context
	parentCtx := context.Background()

	// 1. 设定一个 2 秒的超时期限
	// 我们从 parentCtx 创建一个带有 2 秒超时的子 context
	// Go 语言的惯例是，把 cancel 函数用 defer 调用，确保无论函数如何退出，资源都会被清理
	ctx, cancel := context.WithTimeout(parentCtx, 6*time.Second)
	defer cancel()

	fmt.Println("--- 准备调用 API，超时设置为 2 秒 ---")

	// 2. 将带有超时的 context 传递给我们的函数
	result, err := callSlowAPI(ctx)

	// 3. 检查返回的错误
	if err != nil {
		fmt.Printf("\n调用 API 出错: %v\n", err)

		// 我们可以专门判断错误是不是因为超时引起的
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("错误类型是：任务超时 (context deadline exceeded)")
		}
	} else {
		fmt.Printf("\n成功: %s\n", result)
	}
}
