package user_subscription

import (
	"github.com/tian841224/stock-bot/internal/db/models"
	"github.com/tian841224/stock-bot/internal/repository"
)

// UserSubscriptionService 使用者訂閱服務介面
type UserSubscriptionService interface {
	GetUserSubscriptionByItem(userID uint, item models.SubscriptionItem) (*models.Subscription, error)
	AddUserSubscriptionItem(userID uint, item models.SubscriptionItem) error
	UpdateUserSubscriptionItem(userID uint, item models.SubscriptionItem, status bool) error
	GetUserSubscriptionList(userID uint) ([]*models.Subscription, error)
	AddUserSubscriptionStock(userID uint, stockSymbol string) (bool, error)
	DeleteUserSubscriptionStock(userID uint, stockSymbol string) (bool, error)
	GetUserSubscriptionStockList(userID uint) ([]*repository.UserSubscriptionStock, error)
}

type userSubscriptionService struct {
	userSubscriptionRepo repository.UserSubscriptionRepository
}

func NewUserSubscriptionService(userSubscriptionRepo repository.UserSubscriptionRepository) *userSubscriptionService {
	return &userSubscriptionService{userSubscriptionRepo: userSubscriptionRepo}
}

// GetUserSubscriptionByItem 根據使用者 ID 和訂閱項目取得訂閱資料
func (s *userSubscriptionService) GetUserSubscriptionByItem(userID uint, item models.SubscriptionItem) (*models.Subscription, error) {
	return s.userSubscriptionRepo.GetUserSubscriptionByItem(userID, item)
}

// AddUserSubscriptionItem 新增使用者訂閱項目
func (s *userSubscriptionService) AddUserSubscriptionItem(userID uint, item models.SubscriptionItem) error {
	return s.userSubscriptionRepo.AddUserSubscriptionItem(userID, item)
}

// UpdateUserSubscriptionItem 更新使用者訂閱項目狀態
func (s *userSubscriptionService) UpdateUserSubscriptionItem(userID uint, item models.SubscriptionItem, status bool) error {
	return s.userSubscriptionRepo.UpdateUserSubscriptionItem(userID, item, status)
}

// GetUserSubscriptionList 取得使用者訂閱項目列表
func (s *userSubscriptionService) GetUserSubscriptionList(userID uint) ([]*models.Subscription, error) {
	return s.userSubscriptionRepo.GetUserSubscriptionList(userID)
}

// AddUserSubscriptionStock 新增使用者訂閱股票
func (s *userSubscriptionService) AddUserSubscriptionStock(userID uint, stockSymbol string) (bool, error) {
	return s.userSubscriptionRepo.AddUserSubscriptionStock(userID, stockSymbol)
}

// DeleteUserSubscriptionStock 刪除使用者訂閱股票
func (s *userSubscriptionService) DeleteUserSubscriptionStock(userID uint, stockSymbol string) (bool, error) {
	return s.userSubscriptionRepo.DeleteUserSubscriptionStock(userID, stockSymbol)
}

// GetUserSubscriptionStockList 取得使用者訂閱股票列表
func (s *userSubscriptionService) GetUserSubscriptionStockList(userID uint) ([]*repository.UserSubscriptionStock, error) {
	return s.userSubscriptionRepo.GetUserSubscriptionStockList(userID)
}
