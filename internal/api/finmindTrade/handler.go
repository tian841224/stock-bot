package finmindTrade

import (
	"net/http"
	"stock-bot/internal/infrastructure/finmindtrade"
	"stock-bot/internal/infrastructure/finmindtrade/dto"

	"github.com/gin-gonic/gin"
)

// FinmindTradeHandler FinMind Trade API處理器
type FinmindTradeHandler struct {
	finmindClient finmindtrade.FinmindTradeAPIInterface
}

// NewFinmindTradeHandler 建立新的FinMind Trade處理器
func NewFinmindTradeHandler(finmindClient finmindtrade.FinmindTradeAPIInterface) *FinmindTradeHandler {
	return &FinmindTradeHandler{
		finmindClient: finmindClient,
	}
}

// GetTaiwanStockInfo 取得台灣股票資訊
func (h *FinmindTradeHandler) GetTaiwanStockInfo(c *gin.Context) {
	result, err := h.finmindClient.GetTaiwanStockInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票資訊失敗",
			"code":   "TAIWAN_STOCK_INFO_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票資訊",
	})
}

// GetTaiwanStockPrice 取得台灣股票價格
func (h *FinmindTradeHandler) GetTaiwanStockPrice(c *gin.Context) {
	requestDto := dto.FinmindtradeRequestDto{
		DataID:    c.Query("data_id"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	result, err := h.finmindClient.GetTaiwanStockPrice(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票價格失敗",
			"code":   "TAIWAN_STOCK_PRICE_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票價格",
	})
}

// GetTaiwanExchangeRate 取得台灣匯率
func (h *FinmindTradeHandler) GetTaiwanExchangeRate(c *gin.Context) {
	stockID := c.Query("stock_id")
	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "stock_id為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanExchangeRate(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣匯率失敗",
			"code":   "TAIWAN_EXCHANGE_RATE_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣匯率",
	})
}

// GetTaiwanStockDividend 取得台灣股票股利
func (h *FinmindTradeHandler) GetTaiwanStockDividend(c *gin.Context) {
	stockID := c.Query("stock_id")
	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "stock_id為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanStockDividend(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票股利失敗",
			"code":   "TAIWAN_STOCK_DIVIDEND_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票股利",
	})
}

// GetTaiwanStockFinancialStatements 取得台灣股票財務報表
func (h *FinmindTradeHandler) GetTaiwanStockFinancialStatements(c *gin.Context) {
	stockID := c.Query("stock_id")
	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "stock_id為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanStockFinancialStatements(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票財務報表失敗",
			"code":   "TAIWAN_STOCK_FINANCIAL_STATEMENTS_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票財務報表",
	})
}

// GetTaiwanStockMonthRevenue 取得台灣股票月營收
func (h *FinmindTradeHandler) GetTaiwanStockMonthRevenue(c *gin.Context) {
	stockID := c.Query("stock_id")
	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "stock_id為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanStockMonthRevenue(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票月營收失敗",
			"code":   "TAIWAN_STOCK_MONTH_REVENUE_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票月營收",
	})
}

// GetTaiwanStockTradingDate 取得台灣股票交易日
func (h *FinmindTradeHandler) GetTaiwanStockTradingDate(c *gin.Context) {
	stockID := c.Query("stock_id")

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanStockTradingDate(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票交易日失敗",
			"code":   "TAIWAN_STOCK_TRADING_DATE_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票交易日",
	})
}

// GetTaiwanVariousIndicators 取得台灣各種指標
func (h *FinmindTradeHandler) GetTaiwanVariousIndicators(c *gin.Context) {
	stockID := c.Query("stock_id")

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanVariousIndicators(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣各種指標失敗",
			"code":   "TAIWAN_VARIOUS_INDICATORS_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣各種指標",
	})
}

// GetUSStockInfo 取得美股股票資訊
func (h *FinmindTradeHandler) GetUSStockInfo(c *gin.Context) {
	result, err := h.finmindClient.GetUSStockInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得美股股票資訊失敗",
			"code":   "US_STOCK_INFO_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得美股股票資訊",
	})
}

// GetUSStockPrice 取得美股股票價格
func (h *FinmindTradeHandler) GetUSStockPrice(c *gin.Context) {
	stockID := c.Query("stock_id")

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetUSStockPrice(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得美股股票價格失敗",
			"code":   "US_STOCK_PRICE_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得美股股票價格",
	})
}

// GetTodayInfo 取得今日資訊
func (h *FinmindTradeHandler) GetTodayInfo(c *gin.Context) {
	result, err := h.finmindClient.GetTodayInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得今日資訊失敗",
			"code":   "TODAY_INFO_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得今日資訊",
	})
}

// GetTaiwanStockAnalysis 取得台灣股票分析
func (h *FinmindTradeHandler) GetTaiwanStockAnalysis(c *gin.Context) {
	stockID := c.Query("stock_id")

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanStockAnalysis(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票分析失敗",
			"code":   "TAIWAN_STOCK_ANALYSIS_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票分析",
	})
}

// GetTaiwanStockAnalysisPlot 取得台灣股票分析圖表
func (h *FinmindTradeHandler) GetTaiwanStockAnalysisPlot(c *gin.Context) {
	stockID := c.Query("stock_id")

	requestDto := dto.FinmindtradeRequestDto{
		StockID: stockID,
	}

	result, err := h.finmindClient.GetTaiwanStockAnalysisPlot(requestDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得台灣股票分析圖表失敗",
			"code":   "TAIWAN_STOCK_ANALYSIS_PLOT_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"message": "成功取得台灣股票分析圖表",
	})
}

// GetHealthCheck 健康檢查
func (h *FinmindTradeHandler) GetHealthCheck(c *gin.Context) {
	// 嘗試取得台灣股票資訊來測試API連線
	_, err := h.finmindClient.GetTaiwanStockInfo()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": "finmind-trade-api",
			"error":   err.Error(),
			"message": "FinMind Trade API服務異常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "finmind-trade-api",
		"message": "FinMind Trade API服務正常",
	})
}
