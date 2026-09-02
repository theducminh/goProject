package postgres

import (
	"context"

	"goproject/domain"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repository *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return repository.db.WithContext(ctx).Create(product).Error
}

func (repository *ProductRepository) GetByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	if err := repository.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (repository *ProductRepository) GetByCode(ctx context.Context, code string) (*domain.Product, error) {
	var product domain.Product
	if err := repository.db.WithContext(ctx).Where("code = ?", code).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (repository *ProductRepository) List(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := repository.db.WithContext(ctx).Order("id DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repository *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return repository.db.WithContext(ctx).Save(product).Error
}

func (repository *ProductRepository) Delete(ctx context.Context, id uint) error {
	result := repository.db.WithContext(ctx).Delete(&domain.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (repository *ProductRepository) Ping(ctx context.Context) error {
	database, err := repository.db.DB()
	if err != nil {
		return err
	}
	return database.PingContext(ctx)
}
