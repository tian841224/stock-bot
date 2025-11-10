package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tian841224/stock-bot/config"
	"github.com/tian841224/stock-bot/internal/api/linebot"
	"github.com/tian841224/stock-bot/internal/api/tgbot"
	"github.com/tian841224/stock-bot/internal/db"
	cnyesInfra "github.com/tian841224/stock-bot/internal/infrastructure/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade"
	fugleInfra "github.com/tian841224/stock-bot/internal/infrastructure/fugle"
	"github.com/tian841224/stock-bot/internal/infrastructure/imgbb"
	linebotInfra "github.com/tian841224/stock-bot/internal/infrastructure/linebot"
	tgbotInfra "github.com/tian841224/stock-bot/internal/infrastructure/tgbot"
	twseInfra "github.com/tian841224/stock-bot/internal/infrastructure/twse"
	"github.com/tian841224/stock-bot/internal/repository"
	lineService "github.com/tian841224/stock-bot/internal/service/bot/line"
	tgService "github.com/tian841224/stock-bot/internal/service/bot/tg"
	twstockService "github.com/tian841224/stock-bot/internal/service/twstock"
	"github.com/tian841224/stock-bot/internal/service/user"
	"github.com/tian841224/stock-bot/internal/service/user_subscription"
	"github.com/tian841224/stock-bot/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 初始化結果結構
type InitResult struct {
	cfg                  *config.Config
	log                  logger.Logger
	userRepo             repository.UserRepository
	symbolsRepo          repository.SymbolRepository
	userSubscriptionRepo repository.UserSubscriptionRepository
	fugleAPI             *fugleInfra.FugleAPI
	finmindClient        *finmindtrade.FinmindTradeAPI
	twseAPI              *twseInfra.TwseAPI
	cnyesAPI             *cnyesInfra.CnyesAPI
	imgbbClient          *imgbb.ImgBBClient
	userService          user.UserService
	stockService         twstockService.StockService
	lineBotClient        *linebotInfra.LineBotClient
	tgBotClient          *tgbotInfra.TgBotClient
	err                  error
}

func main() {
	// 初始化日誌
	log, err := logger.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("初始化日誌失敗: %v", err))
	}
	defer log.Sync()

	// 非同步初始化
	initResult, err := asyncInit(log)
	if err != nil {
		log.Panic("初始化失敗", zap.Error(err))
	}

	// 設定 Gin 模式（根據環境變數自動設定）
	// 在 Docker 環境中，GIN_MODE 環境變數會自動設定為 release
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // 預設為 debug 模式（開發環境）
	}
	gin.SetMode(ginMode)
	log.Info("Gin 模式設定", zap.String("mode", ginMode))

	// 建立 Gin 引擎與註冊路由
	router := gin.Default()

	// 健康檢查
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	// 建立使用者訂閱服務
	userSubscriptionService := user_subscription.NewUserSubscriptionService(initResult.userSubscriptionRepo)
	// 建立 LINE Bot 服務層
	lineSvc := lineService.NewLineService(initResult.stockService, userSubscriptionService, initResult.log)
	lineCommandHandler := lineService.NewLineCommandHandler(
		initResult.lineBotClient,
		lineSvc,
		initResult.userService,
		userSubscriptionService,
		initResult.imgbbClient,
		initResult.log,
	)
	service := lineService.NewBotService(initResult.lineBotClient, lineCommandHandler, initResult.userService, initResult.log)
	handler := linebot.NewLineBotHandler(service, initResult.lineBotClient, initResult.log)
	linebot.RegisterRoutes(router, handler, initResult.cfg.LINE_BOT_WEBHOOK_PATH)

	// 建立 Telegram Bot 服務層
	tgSvc := tgService.NewTgService(initResult.stockService, userSubscriptionService, initResult.log)
	tgCommandHandler := tgService.NewTgCommandHandler(
		initResult.tgBotClient,
		tgSvc,
		initResult.userService,
		userSubscriptionService,
		initResult.log,
	)
	tgServiceHandler := tgService.NewTgServiceHandler(tgCommandHandler, initResult.userService, initResult.log)
	tgHandler := tgbot.NewTgHandler(initResult.cfg, tgServiceHandler, initResult.log)
	tgbot.RegisterRoutes(router, tgHandler, initResult.cfg.TELEGRAM_BOT_WEBHOOK_PATH)

	// 從環境變數讀取埠號，預設 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 啟動伺服器（背景）
	serverErr := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	initResult.log.Info("HTTP 伺服器啟動成功")
	initResult.log.Info("程式執行中...")

	// 等待終止訊號或啟動錯誤
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		// 繼續往下優雅關閉
	case err := <-serverErr:
		initResult.log.Error("啟動 HTTP 伺服器失敗", zap.Error(err))
		// 立刻以非 0 退出，先同步日誌
		_ = initResult.log.Sync()
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(ctx)
	if shutdownErr != nil {
		initResult.log.Error("伺服器關閉失敗", zap.Error(shutdownErr))
	}

	dbErr := db.Close()
	if dbErr != nil {
		initResult.log.Error("資料庫關閉失敗", zap.Error(dbErr))
	}

	if shutdownErr != nil || dbErr != nil {
		// 有關閉錯誤，用非 0 退出；先同步日誌
		_ = initResult.log.Sync()
		os.Exit(1)
	}
}

// 非同步初始化函數
func asyncInit(log logger.Logger) (*InitResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := &InitResult{log: log}
	var wg sync.WaitGroup

	// 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("載入設定失敗: %v", err)
	}
	result.cfg = cfg
	log.Info("設定載入成功")

	// 初始化資料庫
	if err := db.InitDB(cfg); err != nil {
		return nil, fmt.Errorf("資料庫初始化失敗: %v", err)
	}
	log.Info("資料庫初始化成功")

	// 並行初始化 Repository
	wg.Add(3)
	go func() {
		defer wg.Done()
		result.userRepo = repository.NewUserRepository(db.GetDB())
		log.Info("UserRepository 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.symbolsRepo = repository.NewSymbolRepository(db.GetDB())
		log.Info("SymbolRepository 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.userSubscriptionRepo = repository.NewUserSubscriptionRepository(db.GetDB())
		log.Info("UserSubscriptionRepository 初始化完成")
	}()

	// 並行初始化外部 API 客戶端
	wg.Add(4)
	go func() {
		defer wg.Done()
		result.fugleAPI = fugleInfra.NewFugleAPI(*cfg)
		log.Info("FugleAPI 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.finmindClient = finmindtrade.NewFinmindTradeAPI(*cfg)
		log.Info("FinmindTradeAPI 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.twseAPI = twseInfra.NewTwseAPI()
		log.Info("TwseAPI 初始化完成")
	}()

	go func() {
		defer wg.Done()
		result.cnyesAPI = cnyesInfra.NewCnyesAPI()
		log.Info("CnyesAPI 初始化完成")
	}()

	// 初始化 ImgBB 客戶端（條件性）
	wg.Add(1)
	go func() {
		defer wg.Done()
		if cfg.IMGBB_API_KEY != "" {
			result.imgbbClient = imgbb.NewImgBBClient(cfg.IMGBB_API_KEY)
			log.Info("ImgBB 客戶端初始化成功")
		} else {
			log.Warn("IMGBB_API_KEY 未設定，圖片上傳功能將不可用")
		}
	}()

	// 並行初始化 Bot 客戶端
	wg.Add(2)
	go func() {
		defer wg.Done()
		botClient, err := linebotInfra.NewBot(*cfg, log)
		if err != nil {
			result.err = fmt.Errorf("初始化 LINE Bot 失敗: %v", err)
			return
		}
		result.lineBotClient = botClient
		log.Info("LINE Bot 客戶端初始化完成")
	}()

	go func() {
		defer wg.Done()
		botClient, err := tgbotInfra.NewBot(*cfg, log)
		if err != nil {
			result.err = fmt.Errorf("初始化 Telegram Bot 失敗: %v", err)
			return
		}
		result.tgBotClient = botClient
		log.Info("Telegram Bot 客戶端初始化完成")
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
	result.userService = user.NewUserService(result.userRepo)
	result.stockService = twstockService.NewStockService(
		result.finmindClient,
		result.twseAPI,
		result.cnyesAPI,
		result.fugleAPI,
		result.symbolsRepo,
		log,
	)

	log.Info("所有初始化完成")
	return result, nil
}
