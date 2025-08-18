package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// 模拟发送邮件的函数
func sendEmailSync(to string) {
	fmt.Printf("[同步] 开始发送邮件给 %s...\n", to)
	time.Sleep(2 * time.Second) // 模拟邮件发送耗时
	fmt.Printf("[同步] 邮件发送完成给 %s\n", to)
}

func sendEmailAsync(to string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[异步] 开始发送邮件给 %s...\n", to)
	time.Sleep(2 * time.Second) // 模拟邮件发送耗时
	fmt.Printf("[异步] 邮件发送完成给 %s\n", to)
}

// 模拟留言处理
func handleMessageSync(message string) {
	start := time.Now()

	// 1. 保存留言到数据库
	fmt.Println("保存留言:", message)
	time.Sleep(100 * time.Millisecond) // 模拟数据库操作

	// 2. 发送邮件通知（同步）
	sendEmailSync("admin@example.com")

	// 3. 返回响应
	fmt.Printf("处理完成，总耗时: %v\n", time.Since(start))
}

func handleMessageAsync(message string) {
	start := time.Now()

	// 1. 保存留言到数据库
	fmt.Println("保存留言:", message)
	time.Sleep(100 * time.Millisecond) // 模拟数据库操作

	// 2. 发送邮件通知（异步）
	var wg sync.WaitGroup
	wg.Add(1)
	go sendEmailAsync("admin@example.com", &wg)

	// 3. 立即返回响应（不等待邮件发送）
	fmt.Printf("响应返回，耗时: %v\n", time.Since(start))

	// 在实际应用中，这里不需要等待
	// 但为了演示，我们等待协程完成
	wg.Wait()
	fmt.Printf("邮件发送完成，总耗时: %v\n", time.Since(start))
}

// 演示协程池处理批量任务
func processBatchMessages(messages []string) {
	start := time.Now()

	// 创建一个有缓冲的channel作为协程池
	workerPool := make(chan struct{}, 3) // 最多3个并发
	var wg sync.WaitGroup

	for i, msg := range messages {
		wg.Add(1)
		workerPool <- struct{}{} // 获取一个工作位

		go func(id int, message string) {
			defer func() {
				<-workerPool // 释放工作位
				wg.Done()
			}()

			// 模拟处理留言
			fmt.Printf("协程 %d: 开始处理 - %s\n", id, message)
			time.Sleep(1 * time.Second)
			fmt.Printf("协程 %d: 处理完成 - %s\n", id, message)
		}(i, msg)
	}

	wg.Wait()
	fmt.Printf("\n批量处理完成，总耗时: %v\n", time.Since(start))
}

func main() {
	fmt.Println("=== 留言功能协程演示 ===\n")

	// 1. 同步处理方式
	fmt.Println("1. 同步处理留言:")
	handleMessageSync("这是一条测试留言")

	fmt.Print("\n" + strings.Repeat("-", 50) + "\n\n")

	// 2. 异步处理方式
	fmt.Println("2. 异步处理留言:")
	handleMessageAsync("这是一条测试留言")

	fmt.Print("\n" + strings.Repeat("-", 50) + "\n\n")

	// 3. 批量处理演示
	fmt.Println("3. 批量处理留言（使用协程池）:")
	messages := []string{
		"留言1: 你的博客很棒！",
		"留言2: 请问如何学习Go语言？",
		"留言3: 期待更多文章",
		"留言4: 这个功能很实用",
		"留言5: 感谢分享",
	}
	processBatchMessages(messages)

	fmt.Println("\n=== 演示结束 ===")
}

/*
运行这个演示，你会看到：

1. 同步处理：用户需要等待邮件发送完成（2秒+）才能收到响应
2. 异步处理：用户立即收到响应（约100ms），邮件在后台发送
3. 批量处理：5条留言用协程池并发处理，总时间约2秒（而不是5秒）

这就是协程在留言功能中的实际应用效果！
*/
