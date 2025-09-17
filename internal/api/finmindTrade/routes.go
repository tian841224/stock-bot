package finmindTrade

import (
	"net/http"
	"stock-bot/internal/infrastructure/finmindtrade"
	"stock-bot/internal/infrastructure/finmindtrade/dto"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 註冊 FinMind Trade API 的路由
func RegisterRoutes(r *gin.Engine, finmindClient finmindtrade.FinmindTradeAPIInterface) {
	r.GET("/finmind/taiwan_stock_info", func(c *gin.Context) {
		result, err := finmindClient.GetTaiwanStockInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_price", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			DataID:    c.Query("data_id"),
			StartDate: c.Query("start_date"),
			EndDate:   c.Query("end_date"),
		}
		result, err := finmindClient.GetTaiwanStockPrice(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_exchange_rate", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanExchangeRate(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_dividend", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanStockDividend(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_financial_statements", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanStockFinancialStatements(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_month_revenue", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanStockMonthRevenue(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_trading_date", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanStockTradingDate(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_various_indicators", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanVariousIndicators(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/us_stock_info", func(c *gin.Context) {
		result, err := finmindClient.GetUSStockInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/us_stock_price", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetUSStockPrice(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/today_info", func(c *gin.Context) {
		result, err := finmindClient.GetTodayInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_analysis", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanStockAnalysis(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.GET("/finmind/taiwan_stock_analysis_plot", func(c *gin.Context) {
		requestDto := dto.FinmindtradeRequestDto{
			StockID: c.Query("stock_id"),
		}
		result, err := finmindClient.GetTaiwanStockAnalysisPlot(requestDto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	})
}
