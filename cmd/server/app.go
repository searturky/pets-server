package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pets-server/internal/infrastructure/cron"
	ws "pets-server/internal/interfaces/websocket"
	"pets-server/internal/pkg/config"
)

// App 应用程序结构体，统一管理生命周期
type App struct {
	cfg       *config.Config
	server    *http.Server
	router    *gin.Engine
	wsHub     *ws.Hub
	scheduler *cron.Scheduler
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	router *gin.Engine,
	wsHub *ws.Hub,
	scheduler *cron.Scheduler,
) *App {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &App{
		cfg:       cfg,
		server:    server,
		router:    router,
		wsHub:     wsHub,
		scheduler: scheduler,
	}
}

// Run 启动应用
func (a *App) Run() error {
	// 启动 WebSocket Hub
	go a.wsHub.Run()

	// 启动定时任务
	a.scheduler.Start()

	// 启动 HTTP 服务
	log.Printf("Server starting on %s", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	// 停止定时任务
	a.scheduler.Stop()

	// 优雅关闭 HTTP 服务器
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}
