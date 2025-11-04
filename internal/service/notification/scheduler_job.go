// Package notification 提供通知排程服務
package notification

import (
	"strconv"
	"time"

	tgbotInfra "github.com/tian841224/stock-bot/internal/infrastructure/tgbot"
	"github.com/tian841224/stock-bot/internal/repository"
	tgbot "github.com/tian841224/stock-bot/internal/service/bot/tg"
	"github.com/tian841224/stock-bot/pkg/logger"
	"go.uber.org/zap"
)

type schedulerJobService struct {
	tgService              *tgbot.TgService
	tgClient               *tgbotInfra.TgBotClient
	userRepo               repository.UserRepository
	subscriptionRepo       repository.SubscriptionRepository
	subscriptionSymbolRepo repository.SubscriptionSymbolRepository
}

func NewSchedulerJobService(tgService *tgbot.TgService, tgClient *tgbotInfra.TgBotClient, userRepo repository.UserRepository, subscriptionRepo repository.SubscriptionRepository, subscriptionSymbolRepo repository.SubscriptionSymbolRepository) *schedulerJobService {
	return &schedulerJobService{
		tgService:              tgService,
		tgClient:               tgClient,
		userRepo:               userRepo,
		subscriptionRepo:       subscriptionRepo,
		subscriptionSymbolRepo: subscriptionSymbolRepo,
	}
}

// NotificationStockPrice 通知當日股價資訊
func (s *schedulerJobService) NotificationStockPrice() {
	// 取得按 symbol 分組的訂閱者清單
	symbolSubscriptions, err := s.getSymbolSubscriptions()
	if err != nil {
		return
	}

	if symbolSubscriptions == nil {
		return
	}

	totalSubscriptions := 0
	// 對每個唯一的 symbol 只查一次股票資訊
	for symbol, userIDs := range symbolSubscriptions {
		tgStockInfoMessage, err := s.tgService.GetStockPriceByDate(symbol, time.Now().Format("2006-01-02"))
		if err != nil {
			logger.Log.Error("取得股票資訊失敗", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		// 將股票資訊發送給所有訂閱該 symbol 的使用者
		s.sendNotificationToSubscribers(tgStockInfoMessage, userIDs)
		totalSubscriptions += len(userIDs)
	}

	logger.Log.Info("股票資訊通知完成", zap.Int("symbol數量", len(symbolSubscriptions)), zap.Int("訂閱數量", totalSubscriptions))
}

// NotificationStockNews 通知股票新聞
func (s *schedulerJobService) NotificationStockNews() {
	// 取得按 symbol 分組的訂閱者清單
	symbolSubscriptions, err := s.getSymbolSubscriptions()
	if err != nil {
		return
	}

	if symbolSubscriptions == nil {
		return
	}

	totalSubscriptions := 0
	// 對每個唯一的 symbol 只查一次股票新聞
	for symbol, userIDs := range symbolSubscriptions {
		tgStockNewsMessage, err := s.tgService.GetTaiwanStockNews(symbol)
		if err != nil {
			logger.Log.Error("取得股票新聞失敗", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		// 將股票資訊發送給所有訂閱該 symbol 的使用者
		for _, userID := range userIDs {
			user, err := s.userRepo.GetByID(userID)
			if err != nil {
				logger.Log.Error("取得使用者失敗", zap.Uint("userID", userID), zap.Error(err))
				continue
			}
			if user == nil {
				logger.Log.Error("使用者資料為空", zap.Uint("userID", userID))
				continue
			}
			accountIDInt, err := strconv.ParseInt(user.AccountID, 10, 64)
			if err != nil {
				logger.Log.Error("轉換使用者 AccountID 失敗", zap.String("accountID", user.AccountID), zap.Error(err))
				continue
			}
			if err := s.tgClient.SendMessageWithKeyboard(accountIDInt, tgStockNewsMessage.Text, tgStockNewsMessage.InlineKeyboardMarkup); err != nil {
				logger.Log.Error("發送股票新聞通知失敗", zap.String("symbol", symbol), zap.Uint("userID", userID), zap.Error(err))
			}
		}

		// 將股票新聞發送給所有訂閱該 symbol 的使用者
		// s.sendNotificationToSubscribers(tgStockNewsMessage.Text, userIDs)
		totalSubscriptions += len(userIDs)
	}

	logger.Log.Info("股票新聞通知完成", zap.Int("symbol數量", len(symbolSubscriptions)), zap.Int("訂閱數量", totalSubscriptions))
}

// NotificationDailyMarketInfo 通知大盤資訊
func (s *schedulerJobService) NotificationDailyMarketInfo() {
	// 取得大盤訂閱者清單
	subscriptionsList, err := s.getSubscriptions(3)
	if err != nil {
		return
	}

	if subscriptionsList == nil {
		return
	}
	tgDailyMarketInfoMessage, err := s.tgService.GetDailyMarketInfo(1)

	// 將大盤資訊發送給所有訂閱者
	s.sendNotificationToSubscribers(tgDailyMarketInfoMessage, subscriptionsList)

	logger.Log.Info("大盤資訊通知完成", zap.Int("訂閱數量", len(subscriptionsList)))
}

// NotificationTopVolumeItems 通知當日交易量前20名資訊
func (s *schedulerJobService) NotificationTopVolumeItems() {
	// 取得大盤訂閱者清單
	subscriptionsList, err := s.getSubscriptions(4)
	if err != nil {
		return
	}

	if subscriptionsList == nil {
		return
	}
	tgTopVolumeItemsMessage, err := s.tgService.GetTopVolumeItemsFormatted()

	// 將大盤資訊發送給所有訂閱者
	s.sendNotificationToSubscribers(tgTopVolumeItemsMessage, subscriptionsList)

	logger.Log.Info("交易量前20名資訊通知完成", zap.Int("訂閱數量", len(subscriptionsList)))
}

// getSubscriptions 取得按 symbol 分組的訂閱者清單
func (s *schedulerJobService) getSubscriptions(featureID uint) ([]uint, error) {
	// 取得所有股票訂閱清單
	subscriptionList, err := s.subscriptionRepo.GetByFeatureID(featureID)
	if err != nil {
		logger.Log.Error("取得所有訂閱清單失敗", zap.Error(err))
		return nil, err
	}

	if len(subscriptionList) == 0 {
		logger.Log.Info("沒有資料需要通知")
		return nil, nil
	}

	// 按 symbol 分組，彙整相同 symbol 的所有訂閱
	subscriptions := make([]uint, 0)
	for _, subscription := range subscriptionList {
		userID := subscription.UserID
		subscriptions = append(subscriptions, userID)
	}

	return subscriptions, nil
}

// getSymbolSubscriptions 取得按 symbol 分組的訂閱者清單
func (s *schedulerJobService) getSymbolSubscriptions() (map[string][]uint, error) {
	// 取得所有股票訂閱清單
	subscriptionSymbols, err := s.subscriptionSymbolRepo.GetAll("subscription_id")
	if err != nil {
		logger.Log.Error("取得所有股票訂閱清單失敗", zap.Error(err))
		return nil, err
	}

	if len(subscriptionSymbols) == 0 {
		logger.Log.Info("沒有訂閱資料需要通知")
		return nil, nil
	}

	// 按 symbol 分組，彙整相同 symbol 的所有訂閱
	symbolSubscriptions := make(map[string][]uint)
	for _, subscriptionSymbol := range subscriptionSymbols {
		if subscriptionSymbol.Symbol == nil || subscriptionSymbol.Subscription == nil {
			logger.Log.Warn("訂閱資料缺少關聯資訊，跳過")
			continue
		}
		symbol := subscriptionSymbol.Symbol.Symbol
		userID := subscriptionSymbol.Subscription.UserID
		symbolSubscriptions[symbol] = append(symbolSubscriptions[symbol], userID)
	}

	return symbolSubscriptions, nil
}

// sendNotificationToSubscribers 將訊息發送給指定 symbol 的所有訂閱者
func (s *schedulerJobService) sendNotificationToSubscribers(message string, userIDs []uint) {
	// 將股票資訊發送給所有訂閱該 symbol 的使用者
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			logger.Log.Error("取得使用者失敗", zap.Uint("userID", userID), zap.Error(err))
			continue
		}
		if user == nil {
			logger.Log.Error("使用者資料為空", zap.Uint("userID", userID))
			continue
		}
		accountIDInt, err := strconv.ParseInt(user.AccountID, 10, 64)
		if err != nil {
			logger.Log.Error("轉換使用者 AccountID 失敗", zap.String("accountID", user.AccountID), zap.Error(err))
			continue
		}
		if err := s.tgClient.SendMessage(accountIDInt, message); err != nil {
			logger.Log.Error("發送通知失敗", zap.Uint("userID", userID), zap.Error(err))
		}
	}
}
