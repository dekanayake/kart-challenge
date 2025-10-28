package routes

import (
	"github.com/dekanayake/kart-challenge/backend-challenge/internal"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/controllers"
	"github.com/dekanayake/kart-challenge/backend-challenge/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(server internal.Server) *gin.Engine {
	r := gin.New()
	r.Use(middleware.ZerologMiddleware())
	r.Use(gin.Recovery())

	productController := controllers.NewProductController(*server.ProductRepo)
	orderController := controllers.NewOrderController(*server.OrderRepo, *server.ProductRepo, *server.FileReader)

	api := r.Group("/api")
	{
		api.GET("/health", controllers.HealthHandler)
		api.GET("/product/:productId", productController.GetProductByID)
		api.GET("/product", productController.ListProducts)
		api.POST("/order", orderController.CreateOrder)
	}

	return r
}
