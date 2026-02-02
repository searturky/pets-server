// Package main 程序入口
// 负责依赖注入和服务启动
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化应用（由 Wire 生成）
	app, cleanup, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	// 在 goroutine 中启动应用
	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("App run failed: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭
	if err := app.Shutdown(context.Background()); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}
