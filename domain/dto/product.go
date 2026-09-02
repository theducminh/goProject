package dto

type CreateProductRequest struct {
	Code  string `json:"code" validate:"required,max=64"`
	Name  string `json:"name" validate:"required,max=255"`
	Price uint   `json:"price" validate:"gte=0"`
	Stock uint   `json:"stock" validate:"gte=0"`
}

type UpdateProductRequest struct {
	Code  *string `json:"code" validate:"omitempty,max=64"`
	Name  *string `json:"name" validate:"omitempty,max=255"`
	Price *uint   `json:"price" validate:"omitempty,gte=0"`
	Stock *uint   `json:"stock" validate:"omitempty,gte=0"`
}

type ProductResponse struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
	Stock uint   `json:"stock"`
}
