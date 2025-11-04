package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tian841224/stock-bot/config"
	"github.com/tian841224/stock-bot/internal/db"
	cnyesInfra "github.com/tian841224/stock-bot/internal/infrastructure/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade"
	fugleInfra "github.com/tian841224/stock-bot/internal/infrastructure/fugle"
	tgbotInfra "github.com/tian841224/stock-bot/internal/infrastructure/tgbot"
	twseInfra "github.com/tian841224/stock-bot/internal/infrastructure/twse"
	"github.com/tian841224/stock-bot/internal/repository"
	tgService "github.com/tian841224/stock-bot/internal/service/bot/tg"
	"github.com/tian841224/stock-bot/internal/service/notification"
	twstockService "github.com/tian841224/stock-bot/internal/service/twstock"
	"github.com/tian841224/stock-bot/pkg/logger"

	"go.uber.org/zap"
)

// 初始化結果結構
type InitResult struct {
	cfg                    *config.Config
	symbolsRepo            repository.SymbolRepository
	userRepo               repository.UserRepository
	subscriptionRepo       repository.SubscriptionRepository
	userSubscriptionRepo   repository.UserSubscriptionRepository
	subscriptionSymbolRepo repository.SubscriptionSymbolRepository
	fugleAPI               *fugleInfra.FugleAPI
	finmindClient          *finmindtrade.FinmindTradeAPI
	twseAPI                *twseInfra.TwseAPI
	cnyesAPI               *cnyesInfra.CnyesAPI
	stockService           *twstockService.StockService
	tgBotClient            *tgbotInfra.TgBotClient
	err                    error
}

func main() {
	// 非同步初始化
	initResult, err := asyncInit()
	if err != nil {
		logger.Log.Panic("初始化失敗", zap.Error(err))
	}

	// 建立 Telegram Bot 服務層
	tgSvc := tgService.NewTgService(initResult.stockService, initResult.userSubscriptionRepo)
	// 建立排程通知服務
	schedulerJobService := notification.NewSchedulerJobService(tgSvc, initResult.tgBotClient, initResult.userRepo, initResult.subscriptionRepo, initResult.subscriptionSymbolRepo)

	// 從設定檔載入時區（預設 Asia/Taipei）
	timezone := initResult.cfg.SCHEDULER_TIMEZONE
	if timezone == "" {
		timezone = "Asia/Taipei"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		logger.Log.Error("載入時區失敗", zap.String("timezone", timezone), zap.Error(err))
		os.Exit(1)
	}
	logger.Log.Info("排程時區設定", zap.String("timezone", timezone))

	// 建立 Cron
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithSeconds(), // 如果需要秒級，可開啟；否則移除這行
	)

	// 程式啟動時先執行一次
	logger.Log.Info("程式啟動時執行一次通知")
	schedulerJobService.NotificationStockPrice()
	schedulerJobService.NotificationStockNews()
	schedulerJobService.NotificationDailyMarketInfo()
	schedulerJobService.NotificationTopVolumeItems()

	// 從設定檔載入排程規格（預設每天 15 點，周一至周五）
	cronSpec := initResult.cfg.SCHEDULER_STOCK_SPEC
	if cronSpec == "" {
		cronSpec = "0 0 15 * * 1-5"
	}
	_, err = c.AddFunc(cronSpec, func() {
		go func() {
			schedulerJobService.NotificationStockPrice()
		}()
		go func() {
			schedulerJobService.NotificationStockNews()
		}()
		go func() {
			schedulerJobService.NotificationDailyMarketInfo()
		}()
		go func() {
			schedulerJobService.NotificationTopVolumeItems()
		}()
	})
	if err != nil {
		logger.Log.Panic("註冊排程失敗", zap.Error(err))
	}

	c.Start()
	logger.Log.Info("排程器啟動完成")

	// 等待終止訊號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	dbErr := db.Close()
	if dbErr != nil {
		logger.Log.Error("資料庫關閉失敗", zap.Error(dbErr))
		_ = logger.Log.Sync()
		os.Exit(1)
	}

}

// 非同步初始化函數
func asyncInit() (*InitResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := &InitResult{}
	var wg sync.WaitGroup

	// 初始化日誌
	logger.InitLogger()
	defer logger.Log.Sync()

	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("載入設定失敗: %v", err)
	}
	result.cfg = cfg
	logger.Log.Info("設定載入成功")

	// 初始化資料庫
	if err := db.InitDB(cfg); err != nil {
		return nil, fmt.Errorf("資料庫初始化失敗: %v", err)
	}
	logger.Log.Info("資料庫初始化成功")

	// 並行初始化 Repository
	wg.Add(5)
	go func() {
		defer wg.Done()
		result.symbolsRepo = repository.NewSymbolRepository(db.GetDB())
		logger.Log.Info("SymbolRepository 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.userRepo = repository.NewUserRepository(db.GetDB())
		logger.Log.Info("UserRepository 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.userSubscriptionRepo = repository.NewUserSubscriptionRepository(db.GetDB())
		logger.Log.Info("UserSubscriptionRepository 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.subscriptionSymbolRepo = repository.NewSubscriptionSymbolRepository(db.GetDB())
		logger.Log.Info("SubscriptionSymbolRepository 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.subscriptionRepo = repository.NewSubscriptionRepository(db.GetDB())
		logger.Log.Info("SubscriptionRepository 初始化完成")
	}()

	// 並行初始化外部 API 客戶端
	wg.Add(4)
	go func() {
		defer wg.Done()
		result.fugleAPI = fugleInfra.NewFugleAPI(*cfg)
		logger.Log.Info("FugleAPI 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.finmindClient = finmindtrade.NewFinmindTradeAPI(*cfg)
		logger.Log.Info("FinmindTradeAPI 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.twseAPI = twseInfra.NewTwseAPI()
		logger.Log.Info("TwseAPI 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.cnyesAPI = cnyesInfra.NewCnyesAPI()
		logger.Log.Info("CnyesAPI 初始化完成")
	}()

	// 初始化 Telegram Bot 客戶端
	wg.Add(1)
	go func() {
		defer wg.Done()
		botClient, err := tgbotInfra.NewBot(*cfg)
		if err != nil {
			result.err = fmt.Errorf("初始化 Telegram Bot 失敗: %v", err)
			return
		}
		result.tgBotClient = botClient
		logger.Log.Info("Telegram Bot 客戶端初始化完成")
	}()

	// 等待所有並行初始化完成
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// 所有初始化完成
	case <-ctx.Done():
		return nil, fmt.Errorf("初始化超時: %v", ctx.Err())
	}

	// 檢查是否有錯誤
	if result.err != nil {
		return nil, result.err
	}

	// 初始化服務（依賴前面的結果）
	result.stockService = twstockService.NewStockService(
		result.finmindClient,
		result.twseAPI,
		result.cnyesAPI,
		result.fugleAPI,
		result.symbolsRepo,
	)

	logger.Log.Info("所有初始化完成")
	return result, nil
}
