package fugle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stock-bot/config"
	"stock-bot/internal/infrastructure/fugle/dto"
	"time"
)

type FugleAPIInterface interface {
	GetStockIntradayQuote(requestDto dto.FugleStockQuoteRequestDto) (dto.FugleStockQuoteResponseDto, error)
}

type FugleAPI struct {
	baseURL    string
	client     *http.Client
	httpHeader http.Header
}

func NewFugleAPI(cfg config.Config) *FugleAPI {
	return &FugleAPI{
		baseURL: "https://api.fugle.tw/marketdata/v1.0/stock/",
		client:  &http.Client{Timeout: 10 * time.Second},
		httpHeader: http.Header{
			"X-API-KEY": []string{cfg.FUGLE_API_KEY},
		},
	}
}

// 取得日內股票即時報價
func (f *FugleAPI) GetStockIntradayQuote(requestDto dto.FugleStockQuoteRequestDto) (dto.FugleStockQuoteResponseDto, error) {
	url := f.baseURL + "/intraday/quote/" + requestDto.Symbol
	if requestDto.Type != "" {
		url += "?type=" + requestDto.Type
	}
	return getResponse[dto.FugleStockQuoteResponseDto](f, url)
}

// 取得盤中 K 線
func (f *FugleAPI) GetStockIntradayCandles(requestDto dto.FugleCandlesRequestDto) (dto.FugleCandlesResponseDto, error) {
	url := f.baseURL + "/intraday/candles/" + requestDto.Symbol
	if requestDto.From != "" {
		url += "?from=" + requestDto.From
	}
	if requestDto.To != "" {
		url += "?to=" + requestDto.To
	}
	if requestDto.Timeframe != "" {
		url += "?timeframe=" + requestDto.Timeframe
	}
	if requestDto.Fields != "" {
		url += "?fields=" + requestDto.Fields
	}
	if requestDto.Sort != "" {
		url += "?sort=" + requestDto.Sort
	}
	return getResponse[dto.FugleCandlesResponseDto](f, url)
}

// 取得股票歷史Ｋ線
func (f *FugleAPI) GetStockHistoricalCandles(requestDto dto.FugleCandlesRequestDto) (dto.FugleCandlesResponseDto, error) {
	url := f.baseURL + "/historical/candles/" + requestDto.Symbol
	if requestDto.From != "" {
		url += "?from=" + requestDto.From
	}
	if requestDto.To != "" {
		url += "?to=" + requestDto.To
	}
	if requestDto.Timeframe != "" {
		url += "?timeframe=" + requestDto.Timeframe
	}
	if requestDto.Fields != "" {
		url += "?fields=" + requestDto.Fields
	}
	if requestDto.Sort != "" {
		url += "?sort=" + requestDto.Sort
	}
	return getResponse[dto.FugleCandlesResponseDto](f, url)
}

// 取得股票漲跌幅排行快照(需開發者權限)
func (f *FugleAPI) GetStockSnapshotMovers(requestDto dto.FugleMoversRequestDto) (dto.FugleMoversResponseDto, error) {
	url := f.baseURL + "/snapshot/movers"
	if requestDto.Market != "" {
		url += "?market=" + requestDto.Market
	}
	if requestDto.Direction != "" {
		url += "?direction=" + requestDto.Direction
	}
	if requestDto.Change != "" {
		url += "?change=" + requestDto.Change
	}
	if requestDto.Type != "" {
		url += "?type=" + requestDto.Type
	}
	if requestDto.Gt != "" {
		url += "?gt=" + requestDto.Gt
	}
	if requestDto.Gte != "" {
		url += "?gte=" + requestDto.Gte
	}
	if requestDto.Lt != "" {
		url += "?lt=" + requestDto.Lt
	}
	if requestDto.Lte != "" {
		url += "?lte=" + requestDto.Lte
	}
	if requestDto.Eq != "" {
		url += "?eq=" + requestDto.Eq
	}
	return getResponse[dto.FugleMoversResponseDto](f, url)
}

func getResponse[T any](c *FugleAPI, url string) (response T, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}

	req.Header = c.httpHeader

	resp, err := c.client.Do(req)

	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("外部 API 回應錯誤，狀態碼: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("無法讀取 API 回應: %v", err)
	}
	fmt.Printf("Fugle API Response: %s\n", string(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return response, fmt.Errorf("無法解析回應 JSON: %v", err)
	}

	return response, nil
}
