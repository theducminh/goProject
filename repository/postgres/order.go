package postgres

import (
	"context"

	"goproject/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repository *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return repository.db.WithContext(ctx).Create(order).Error
}

func (repository *OrderRepository) GetByID(ctx context.Context, id uint) (*domain.Order, error) {
	var order domain.Order
	if err := repository.db.WithContext(ctx).Preload("Items").First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}
