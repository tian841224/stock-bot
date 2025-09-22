package fugle

import (
	twstockService "stock-bot/internal/service/twstock"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 註冊 Fugle API 的路由
func RegisterRoutes(r *gin.Engine, twstockService *twstockService.StockService) {
	// 建立處理器
	handler := NewFugleHandler(twstockService)

	// 建立路由群組
	fugleGroup := r.Group("/fugle")
	{
		// 健康檢查
		fugleGroup.GET("/health", handler.GetHealthCheck)

		// 股票相關路由
		stockGroup := fugleGroup.Group("/stock")
		{
			// 取得盤中資訊
			stockGroup.GET("/:stock_id/intraday/quote", handler.GetIntradayQuote)
			// 取得歷史Ｋ線資訊
			stockGroup.GET("/:stock_id/historical/candles", handler.GetHistoricalCandles)

		}
	}
}
