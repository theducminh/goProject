package main

// @Summary Create a product
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product payload"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func swaggerCreateProduct() {}

// @Summary List products
// @Produce json
// @Success 200 {array} dto.ProductResponse
// @Failure 500 {object} map[string]string
// @Router /products [get]
func swaggerListProducts() {}

// @Summary Get a product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func swaggerGetProduct() {}

// @Summary Update a product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.UpdateProductRequest true "Product update"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func swaggerUpdateProduct() {}

// @Summary Delete a product
// @Produce json
// @Param id path int true "Product ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func swaggerDeleteProduct() {}

// @Summary Search products
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} dto.ProductResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [get]
func swaggerSearchProducts() {}
