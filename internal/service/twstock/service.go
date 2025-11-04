package twstock

import (
	"github.com/tian841224/stock-bot/internal/infrastructure/cnyes"
	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade"
	"github.com/tian841224/stock-bot/internal/infrastructure/finmindtrade/dto"
	"github.com/tian841224/stock-bot/internal/infrastructure/fugle"
	fugleDto "github.com/tian841224/stock-bot/internal/infrastructure/fugle/dto"
	"github.com/tian841224/stock-bot/internal/infrastructure/twse"
	twseDto "github.com/tian841224/stock-bot/internal/infrastructure/twse/dto"
	"github.com/tian841224/stock-bot/internal/repository"
	stockDto "github.com/tian841224/stock-bot/internal/service/twstock/dto"
)

// StockService 股票服務介面
type StockService interface {
	GetStockPrice(stockID string, date ...string) (*stockDto.StockPriceInfo, error)
	GetStockPerformance(stockID string) (*stockDto.StockPerformanceResponseDto, error)
	GetStockPriceHistory(stockID string) ([]stockDto.StockPerformanceData, error)
	GetAfterTradingVolume(symbol, date string) (*twseDto.AfterTradingVolumeResponseDto, error)
	GetStockNews(stockID string) ([]dto.TaiwanNewsResponseData, error)
	GetStockIntradayQuote(dto fugleDto.FugleStockQuoteRequestDto) (*fugleDto.FugleStockQuoteResponseDto, error)
	GetStockHistoricalCandles(dto fugleDto.FugleCandlesRequestDto) (*fugleDto.FugleCandlesResponseDto, error)
	GetStockInfo(stockID string) (*stockDto.StockQuoteInfo, error)
	GetStockQuote(stockID string) (*stockDto.StockQuoteInfo, error)
	GetTopVolumeItems() ([]*stockDto.StockPriceInfo, error)
	GetStockAnalysis(stockID string) ([]byte, string, error)
	ValidateStockID(stockID string) (bool, string, error)
	GetStockRevenue(stockID string) (*stockDto.RevenueDto, error)
	GetDailyMarketInfo(count int) (twseDto.DailyMarketInfoResponseDto, error)
	GetStockPerformanceWithChart(stockID string, chartType string) (*stockDto.StockPerformanceResponseDto, error)
	GetStockHistoricalCandlesChart(dto fugleDto.FugleCandlesRequestDto) ([]byte, string, error)
	GetStockRevenueChart(stockID string) ([]byte, error)
}

// stockService 股票服務
type stockService struct {
	finmindClient finmindtrade.FinmindTradeAPIInterface
	twseAPI       *twse.TwseAPI
	cnyesAPI      *cnyes.CnyesAPI
	fugleClient   *fugle.FugleAPI
	symbolsRepo   repository.SymbolRepository
	domainService *DomainService
}

// NewStockService 建立股票服務實例
func NewStockService(
	finmindClient finmindtrade.FinmindTradeAPIInterface,
	twseAPI *twse.TwseAPI,
	cnyesAPI *cnyes.CnyesAPI,
	fugleClient *fugle.FugleAPI,
	symbolsRepo repository.SymbolRepository,
) StockService {
	return &stockService{
		finmindClient: finmindClient,
		twseAPI:       twseAPI,
		cnyesAPI:      cnyesAPI,
		fugleClient:   fugleClient,
		symbolsRepo:   symbolsRepo,
		domainService: NewDomainService(),
	}
}

// 確保 stockService 實作了 StockService
var _ StockService = (*stockService)(nil)
