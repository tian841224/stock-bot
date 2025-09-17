package twstock

type StockPerformanceResponseDto struct {
	Data []StockPerformanceData `json:"data"`
}

type StockPerformanceData struct {
	StockID     string `json:"stock_id"`
	Period      string `json:"period"`
	PeriodName  string `json:"period_name"`
	Performance string `json:"performance"`
}
