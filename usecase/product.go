package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"goproject/domain"
	"goproject/domain/dto"
	"goproject/repository"
)

type ProductUsecase interface {
	Create(ctx context.Context, request dto.CreateProductRequest) (dto.ProductResponse, error)
	Get(ctx context.Context, id uint) (dto.ProductResponse, error)
	List(ctx context.Context) ([]dto.ProductResponse, error)
	Update(ctx context.Context, id uint, request dto.UpdateProductRequest) (dto.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string) ([]dto.ProductResponse, error)
}

type productUsecase struct {
	products repository.ProductRepository
	cache    repository.ProductCache
	search   repository.ProductSearchRepository
	timeout  time.Duration
}

func NewProductUsecase(products repository.ProductRepository, cache repository.ProductCache, search repository.ProductSearchRepository, timeout time.Duration) ProductUsecase {
	return &productUsecase{products: products, cache: cache, search: search, timeout: timeout}
}

func (service *productUsecase) Create(parent context.Context, request dto.CreateProductRequest) (dto.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	product := &domain.Product{Code: strings.TrimSpace(request.Code), Name: strings.TrimSpace(request.Name), Price: request.Price, Stock: request.Stock}
	if _, err := service.products.GetByCode(ctx, product.Code); err == nil {
		return dto.ProductResponse{}, domain.ErrAlreadyExists
	} else if err != domain.ErrNotFound {
		return dto.ProductResponse{}, err
	}
	if err := service.products.Create(ctx, product); err != nil {
		return dto.ProductResponse{}, err
	}
	service.cacheProduct(ctx, product)
	return toProductResponse(product), nil
}

func (service *productUsecase) Get(parent context.Context, id uint) (dto.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	if service.cache != nil {
		if product, err := service.cache.Get(ctx, id); err == nil {
			return toProductResponse(product), nil
		}
	}
	product, err := service.products.GetByID(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	service.cacheProduct(ctx, product)
	return toProductResponse(product), nil
}

func (service *productUsecase) List(parent context.Context) ([]dto.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	products, err := service.products.List(ctx)
	if err != nil {
		return nil, err
	}
	responses := make([]dto.ProductResponse, 0, len(products))
	for index := range products {
		responses = append(responses, toProductResponse(&products[index]))
	}
	return responses, nil
}

func (service *productUsecase) Update(parent context.Context, id uint, request dto.UpdateProductRequest) (dto.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	product, err := service.products.GetByID(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	if request.Code != nil {
		product.Code = strings.TrimSpace(*request.Code)
	}
	if request.Name != nil {
		product.Name = strings.TrimSpace(*request.Name)
	}
	if request.Price != nil {
		product.Price = *request.Price
	}
	if request.Stock != nil {
		product.Stock = *request.Stock
	}
	if err := service.products.Update(ctx, product); err != nil {
		return dto.ProductResponse{}, err
	}
	service.cacheProduct(ctx, product)
	return toProductResponse(product), nil
}

func (service *productUsecase) Delete(parent context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	if err := service.products.Delete(ctx, id); err != nil {
		return err
	}
	if service.cache != nil {
		_ = service.cache.Delete(ctx, id)
	}
	if service.search != nil {
		_ = service.search.Delete(ctx, id)
	}
	return nil
}

func (service *productUsecase) Search(parent context.Context, query string) ([]dto.ProductResponse, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	if service.search == nil {
		return nil, errors.New("product search is unavailable")
	}
	products, err := service.search.Search(ctx, strings.TrimSpace(query))
	if err != nil {
		return nil, err
	}
	responses := make([]dto.ProductResponse, 0, len(products))
	for index := range products {
		responses = append(responses, toProductResponse(&products[index]))
	}
	return responses, nil
}

func (service *productUsecase) cacheProduct(ctx context.Context, product *domain.Product) {
	if service.cache != nil {
		_ = service.cache.Set(ctx, product)
	}
	if service.search != nil {
		_ = service.search.Index(ctx, product)
	}
}

func toProductResponse(product *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{ID: product.ID, Code: product.Code, Name: product.Name, Price: product.Price, Stock: product.Stock}
}
