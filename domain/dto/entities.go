package dto

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
}

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required,max=64"`
}

type CreateSupplierRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=32"`
}

type CreateCustomerRequest struct {
	Name  string `json:"name" validate:"required,max=255"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
	Phone string `json:"phone" validate:"omitempty,max=32"`
}

type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required,gt=0"`
	Quantity  uint `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID uint                     `json:"customer_id" validate:"required,gt=0"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}
