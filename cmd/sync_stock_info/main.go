package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tian841224/stock-bot/config"
	"github.com/tian841224/stock-bot/internal/db"
	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade"
	"github.com/tian841224/stock-bot/internal/repository"
	"github.com/tian841224/stock-bot/internal/service/stock_sync"
	"github.com/tian841224/stock-bot/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 初始化日誌
	log, err := logger.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("初始化日誌失敗: %v", err))
	}
	defer log.Sync()

	log.Info("=== 股票資料同步程式啟動 ===")

	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panic("載入設定失敗", zap.Error(err))
	}
	log.Info("設定載入成功")

	// 初始化資料庫
	if err := db.InitDB(cfg); err != nil {
		log.Panic("資料庫初始化失敗", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("資料庫關閉失敗", zap.Error(err))
		}
	}()
	log.Info("資料庫初始化成功")

	// 初始化 Repository 和 Service
	symbolsRepo := repository.NewSymbolRepository(db.GetDB())
	finmindClient := finmindtrade.NewFinmindTradeAPI(*cfg)
	stockSyncService := stock_sync.NewStockSyncService(symbolsRepo, finmindClient, log)
	log.Info("服務初始化成功")

	// 建立 context 用於優雅關閉
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 監聽系統中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 啟動背景同步服務
	go runBackgroundSync(ctx, stockSyncService, log)

	// 等待中斷信號
	<-quit
	log.Info("收到關閉信號，正在優雅關閉...")

	// 取消 context，停止背景任務
	cancel()

	log.Info("=== 程式已關閉 ===")
}

// runBackgroundSync 執行背景同步任務
func runBackgroundSync(ctx context.Context, stockSyncService stock_sync.StockSyncService, log logger.Logger) {
	defer func() {
		log.Info("背景同步任務已完全停止")
	}()

	// 程式啟動時立即執行一次同步
	log.Info("執行初始同步...")
	if err := stockSyncService.SyncTaiwanStockInfo(); err != nil {
		log.Error("初始同步失敗", zap.Error(err))
	}
	log.Info("台股同步完成")
	if err := stockSyncService.SyncUSStockInfo(); err != nil {
		log.Error("初始同步失敗", zap.Error(err))
	}
	log.Info("美股同步完成")

	// 顯示同步統計
	if stats, err := stockSyncService.GetSyncStats(); err == nil {
		log.Info("初始同步統計", zap.Any("stats", stats))
	}

	// 建立定時器
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("收到停止信號，背景同步任務正在關閉...")
			return
		case <-ticker.C:
			log.Info("開始定時同步...")
			if err := stockSyncService.SyncTaiwanStockInfo(); err != nil {
				log.Error("初始同步失敗", zap.Error(err))
			}
			log.Info("台股同步完成")
			if err := stockSyncService.SyncUSStockInfo(); err != nil {
				log.Error("初始同步失敗", zap.Error(err))
			}
			log.Info("美股同步完成")

			// 顯示同步統計
			if stats, err := stockSyncService.GetSyncStats(); err == nil {
				log.Info("初始同步統計", zap.Any("stats", stats))
			}
		}
	}
}
