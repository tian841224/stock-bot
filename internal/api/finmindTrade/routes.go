package finmindTrade

import (
	"stock-bot/internal/infrastructure/finmindtrade"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 註冊 FinMind Trade API 的路由
func RegisterRoutes(r *gin.Engine, finmindClient finmindtrade.FinmindTradeAPIInterface) {
	// 建立處理器
	handler := NewFinmindTradeHandler(finmindClient)

	// 建立路由群組
	finmindGroup := r.Group("/finmind")
	{
		// 健康檢查
		finmindGroup.GET("/health", handler.GetHealthCheck)

		// 台股相關路由
		taiwanGroup := finmindGroup.Group("/taiwan")
		{
			// 基本資訊
			taiwanGroup.GET("/stock_info", handler.GetTaiwanStockInfo)
			taiwanGroup.GET("/stock_price", handler.GetTaiwanStockPrice)
			taiwanGroup.GET("/exchange_rate", handler.GetTaiwanExchangeRate)

			// 財務資訊
			taiwanGroup.GET("/stock_dividend", handler.GetTaiwanStockDividend)
			taiwanGroup.GET("/stock_financial_statements", handler.GetTaiwanStockFinancialStatements)
			taiwanGroup.GET("/stock_month_revenue", handler.GetTaiwanStockMonthRevenue)

			// 交易資訊
			taiwanGroup.GET("/stock_trading_date", handler.GetTaiwanStockTradingDate)
			taiwanGroup.GET("/various_indicators", handler.GetTaiwanVariousIndicators)

			// 分析資訊
			taiwanGroup.GET("/stock_analysis", handler.GetTaiwanStockAnalysis)
			taiwanGroup.GET("/stock_analysis_plot", handler.GetTaiwanStockAnalysisPlot)
		}

		// 美股相關路由
		usGroup := finmindGroup.Group("/us")
		{
			usGroup.GET("/stock_info", handler.GetUSStockInfo)
			usGroup.GET("/stock_price", handler.GetUSStockPrice)
		}

		// 綜合資訊
		finmindGroup.GET("/today_info", handler.GetTodayInfo)
	}
}
