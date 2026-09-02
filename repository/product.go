package repository

import (
	"context"
	"errors"

	"goproject/domain"
)

var ErrCacheMiss = errors.New("cache miss")

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	GetByCode(ctx context.Context, code string) (*domain.Product, error)
	List(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uint) error
}

type ProductCache interface {
	Get(ctx context.Context, id uint) (*domain.Product, error)
	Set(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uint) error
}

type ProductSearchRepository interface {
	Index(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string) ([]domain.Product, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	GetByName(ctx context.Context, name string) (*domain.Role, error)
}

type SupplierRepository interface {
	Create(ctx context.Context, supplier *domain.Supplier) error
	GetByID(ctx context.Context, id uint) (*domain.Supplier, error)
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id uint) (*domain.Customer, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uint) (*domain.Order, error)
}
