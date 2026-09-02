package usecase

import (
	"context"
	"time"

	"goproject/repository"
)

type StockUsecase interface {
	Check(ctx context.Context, productID, quantity uint) (bool, uint, error)
}

type stockUsecase struct {
	products repository.ProductRepository
	timeout  time.Duration
}

func NewStockUsecase(products repository.ProductRepository, timeout time.Duration) StockUsecase {
	return &stockUsecase{products: products, timeout: timeout}
}

func (service *stockUsecase) Check(parent context.Context, productID, quantity uint) (bool, uint, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	product, err := service.products.GetByID(ctx, productID)
	if err != nil {
		return false, 0, err
	}
	return product.Stock >= quantity, product.Stock, nil
}
