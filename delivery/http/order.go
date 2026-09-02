package http

import (
	"errors"
	"net/http"

	"goproject/domain"
	"goproject/domain/dto"
	"goproject/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	orders   usecase.OrderUsecase
	validate *validator.Validate
}

func NewOrderHandler(orders usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orders: orders, validate: validator.New()}
}

func (handler *OrderHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/orders", handler.Create)
}

// CreateOrder godoc
// @Summary Create an order
// @Description Creates an order and publishes an OrderCreated event.
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderRequest true "Order payload"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (handler *OrderHandler) Create(context *gin.Context) {
	var request dto.CreateOrderRequest
	if err := context.ShouldBindJSON(&request); err != nil || handler.validate.Struct(request) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid order data"})
		return
	}
	order, err := handler.orders.Create(context.Request.Context(), request)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	context.JSON(http.StatusCreated, order)
}
