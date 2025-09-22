package fugle

import (
	"net/http"
	"stock-bot/config"
	"stock-bot/internal/infrastructure/fugle"
	"stock-bot/internal/infrastructure/fugle/dto"
	"stock-bot/internal/service/twstock"

	"github.com/gin-gonic/gin"
)

// FugleHandler Fugle API 處理器
type FugleHandler struct {
	twstockService *twstock.StockService
}

// NewFugleHandler 建立新的 Fugle 處理器
func NewFugleHandler(twstockService *twstock.StockService) *FugleHandler {
	return &FugleHandler{
		twstockService: twstockService,
	}
}

// GetHealthCheck 健康檢查
// @Summary 健康檢查
// @Description 檢查 Fugle API 服務是否正常運作
// @Tags Fugle API
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服務正常"
// @Failure 503 {object} map[string]interface{} "服務異常"
// @Router /fugle/health [get]
func (h *FugleHandler) GetHealthCheck(c *gin.Context) {
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": "fugle-api",
			"error":   "無法讀取設定",
			"message": "Fugle API 服務異常",
		})
		return
	}
	fugleAPI := fugle.NewFugleAPI(*cfg)
	req := dto.FugleStockQuoteRequestDto{
		Symbol: "2330",
	}
	_, err = fugleAPI.GetStockIntradayQuote(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": "fugle-api",
			"error":   err.Error(),
			"message": "Fugle API 服務異常",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "fugle-api",
		"message": "Fugle API 服務正常",
	})
}

// GetIntradayQuote 取得盤中即時資料
// @Summary 取得股票盤中即時資料
// @Description 透過 Fugle API 取得指定股票的盤中即時資料
// @Tags Fugle API
// @Accept json
// @Produce json
// @Param stock_id path string true "股票代碼 (例如: 2330)"
// @Success 200 {object} map[string]interface{} "成功取得盤中資料"
// @Failure 400 {object} map[string]interface{} "請求參數錯誤"
// @Failure 500 {object} map[string]interface{} "內部伺服器錯誤"
// @Router /fugle/stock/{stock_id}/intraday/quote [get]
func (h *FugleHandler) GetIntradayQuote(c *gin.Context) {
	stockID := c.Param("stock_id")

	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "股票代碼為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "讀取設定失敗",
			"code":   "CONFIG_ERROR",
			"detail": err.Error(),
		})
		return
	}

	fugleAPI := fugle.NewFugleAPI(*cfg)

	req := dto.FugleStockQuoteRequestDto{
		Symbol: stockID,
	}

	quote, err := fugleAPI.GetStockIntradayQuote(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得盤中資料失敗",
			"code":   "INTERNAL_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (h *FugleHandler) GetHistoricalCandles(c *gin.Context) {
	stockID := c.Param("stock_id")

	if stockID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "股票代碼為必填參數",
			"code":  "MISSING_STOCK_ID",
		})
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	timeframe := c.Query("timeframe")
	fields := c.Query("fields")
	sort := c.Query("sort")

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "讀取設定失敗",
			"code":   "CONFIG_ERROR",
			"detail": err.Error(),
		})
		return
	}

	fugleAPI := fugle.NewFugleAPI(*cfg)

	req := dto.FugleCandlesRequestDto{
		Symbol:    stockID,
		From:      from,
		To:        to,
		Timeframe: timeframe,
		Fields:    fields,
		Sort:      sort,
	}

	candles, err := fugleAPI.GetStockHistoricalCandles(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "取得Ｋ線資料失敗",
			"code":   "INTERNAL_ERROR",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, candles)
}
