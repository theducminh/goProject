package routes

import (
	"context"
	"time"

	productHTTP "goproject/delivery/http"
	"goproject/repository"
	"goproject/usecase"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(products usecase.ProductUsecase, orders usecase.OrderUsecase, checks ...repository.HealthChecker) *gin.Engine {
	r := gin.Default()
	r.GET("/health/live", func(context *gin.Context) { context.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/health/ready", func(ginContext *gin.Context) {
		checkContext, cancel := context.WithTimeout(ginContext.Request.Context(), 2*time.Second)
		defer cancel()
		for _, check := range checks {
			if check != nil && check.Ping(checkContext) != nil {
				ginContext.JSON(503, gin.H{"status": "not_ready"})
				return
			}
		}
		ginContext.JSON(200, gin.H{"status": "ready"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	productHTTP.NewProductHandler(products).RegisterRoutes(r)
	productHTTP.NewOrderHandler(orders).RegisterRoutes(r)

	return r
}
