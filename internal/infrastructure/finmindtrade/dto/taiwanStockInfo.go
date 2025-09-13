package dto

// 台股股票清單
type TaiwanStockInfoResponseDto struct {
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
	Data   []StockInfo `json:"data"`
}

type StockInfo struct {
	IndustryCategory string `json:"industry_category"`
	StockID          string `json:"stock_id"`
	StockName        string `json:"stock_name"`
	Type             string `json:"type"`
	Date             string `json:"date"`
}
