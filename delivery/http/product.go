package http

import (
	"errors"
	"net/http"
	"strconv"

	"goproject/domain"
	"goproject/domain/dto"
	"goproject/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	products  usecase.ProductUsecase
	validator *validator.Validate
}

func NewProductHandler(products usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{products: products, validator: validator.New()}
}

func (handler *ProductHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/products", handler.List)
	router.POST("/products", handler.Create)
	router.GET("/products/:id", handler.Get)
	router.PUT("/products/:id", handler.Update)
	router.DELETE("/products/:id", handler.Delete)
	router.GET("/search", handler.Search)
}

// SearchProducts godoc
// @Summary Search products
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [get]
func (handler *ProductHandler) Search(context *gin.Context) {
	query := context.Query("q")
	if err := handler.validator.Var(query, "required,max=255"); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid search query"})
		return
	}
	products, err := handler.products.Search(context.Request.Context(), query)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, products)
}

// CreateProduct godoc
// @Summary Create a product
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product payload"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (handler *ProductHandler) Create(context *gin.Context) {
	var request dto.CreateProductRequest
	if err := context.ShouldBindJSON(&request); err != nil || handler.validator.Struct(request) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product data"})
		return
	}
	product, err := handler.products.Create(context.Request.Context(), request)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusCreated, product)
}

// ListProducts godoc
// @Summary List products
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (handler *ProductHandler) List(context *gin.Context) {
	products, err := handler.products.List(context.Request.Context())
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (handler *ProductHandler) Get(context *gin.Context) {
	id, ok := parseID(context)
	if !ok {
		return
	}
	product, err := handler.products.Get(context.Request.Context(), id)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.UpdateProductRequest true "Product update"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func (handler *ProductHandler) Update(context *gin.Context) {
	id, ok := parseID(context)
	if !ok {
		return
	}
	var request dto.UpdateProductRequest
	if err := context.ShouldBindJSON(&request); err != nil || handler.validator.Struct(request) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product data"})
		return
	}
	product, err := handler.products.Update(context.Request.Context(), id, request)
	if err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Produce json
// @Param id path int true "Product ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func (handler *ProductHandler) Delete(context *gin.Context) {
	id, ok := parseID(context)
	if !ok {
		return
	}
	if err := handler.products.Delete(context.Request.Context(), id); err != nil {
		handler.writeError(context, err)
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func parseID(context *gin.Context) (uint, bool) {
	parsed, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil || parsed == 0 || parsed > uint64(^uint(0)) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return 0, false
	}
	return uint(parsed), true
}

func (handler *ProductHandler) writeError(context *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		context.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
	case errors.Is(err, domain.ErrAlreadyExists):
		context.JSON(http.StatusConflict, gin.H{"error": "product code already exists"})
	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
