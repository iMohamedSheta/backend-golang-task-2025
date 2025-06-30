package routes

import (
	"taskgo/internal/api/handlers"
	"taskgo/internal/api/middleware"
	"taskgo/internal/deps"
	"taskgo/pkg/ioc"
	"taskgo/pkg/response"
	"taskgo/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynqmon"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	// Global middlewares
	r.Use(middleware.RecoveryWithLogger())
	r.Use(middleware.Logger())
	r.Use(middleware.RateLimiter())
	r.Use(middleware.CORSMiddleware())

	// web routes (without auth)
	registerWebRoutes(r)

	api := r.Group("/api/v1")

	websocketApi := api.Group("/ws")
	registerWebsocketRoutes(websocketApi)

	// User Authentication
	authHandler := deps.App[*handlers.AuthHandler]()
	api.POST("/login", middleware.HandleErrors(authHandler.Login))                      // Done
	api.POST("/refresh-token", middleware.HandleErrors(authHandler.RefreshAccessToken)) // Done

	// Create User (register customer)
	userHandler := deps.App[*handlers.UserHandler]()
	api.POST("/users", middleware.HandleErrors(userHandler.CreateUser)) // Done

	// Product view
	productHandler := deps.App[*handlers.ProductHandler]()
	api.GET("/products", middleware.HandleErrors(productHandler.ListProducts))   // Done
	api.GET("/products/:id", middleware.HandleErrors(productHandler.GetProduct)) // Done
	api.GET("/products/:id/inventory", productHandler.CheckInventory)            // Skipped

	// Protected routes with auth middleware
	api.Use(middleware.Auth())
	{
		// Admin Endpoints (prefix: /api/v1/admin)
		adminApi := api.Group("/admin", middleware.AdminOnly())
		{
			// Admin Order Management
			adminOrderHandler := handlers.AdminOrderHandler{}
			adminApi.GET("/orders", adminOrderHandler.ListAllOrders)
			adminApi.PUT("/orders/:id/status", adminOrderHandler.UpdateOrderStatus)
			adminApi.GET("/reports/daily", adminOrderHandler.DailySalesReport)
			adminApi.GET("/inventory/low-stock", adminOrderHandler.LowStockAlerts)

			// Should make inventory management
			// ...
		}

		// User Management
		api.GET("/users/:id", middleware.HandleErrors(userHandler.GetUser))    // Done
		api.PUT("/users/:id", middleware.HandleErrors(userHandler.UpdateUser)) // Done

		// Product Management
		api.POST("/products", middleware.AdminOnly(), middleware.HandleErrors(productHandler.CreateProduct))    // Done
		api.PUT("/products/:id", middleware.AdminOnly(), middleware.HandleErrors(productHandler.UpdateProduct)) // Done

		// Order Management (VIP)
		orderHandler := deps.App[*handlers.OrderHandler]()
		api.POST("/orders", middleware.HandleErrors(orderHandler.CreateOrder)) // Working on it
		api.GET("/orders", orderHandler.ListUserOrders)
		api.GET("/orders/:id", orderHandler.GetOrder)
		api.PUT("/orders/:id/cancel", orderHandler.CancelOrder)
		api.GET("/orders/:id/status", orderHandler.GetOrderStatus)
	}

	return r
}

func registerWebsocketRoutes(wsApi *gin.RouterGroup) {
	// WebSocket Notifications
	notificationHandler := deps.App[*handlers.NotificationHandler]()
	wsApi.Use(middleware.WebSocketAuth())
	{
		wsApi.GET("/notifications", notificationHandler.WsUserNotificationHandler)
	}
}

func registerWebRoutes(web *gin.Engine) {
	// Root health check
	web.GET("/health", func(c *gin.Context) {
		deps.Log().Log().Info("Health check hit")
		response.Json(c, "Welcome to the Order Processing System", gin.H{"status": "healthy"}, 200)
	})

	// If it's not in production, show swagger docs
	if deps.Config().GetString("app.env", "prod") != "prod" {
		web.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	}

	// Monitoring asynq jobs
	ops, _ := utils.GetRedisQueueClientOptionsForAsynq("redis.connections.queue")
	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitor",
		RedisConnOpt: ops,
	})

	web.GET("/monitor/*a", gin.WrapH(h))

	web.GET("/dependencies", func(c *gin.Context) {
		deps := ioc.AppContainer().GetRegisteredService()
		response.Json(c, "Dependencies loaded", deps, 200)
	})
	notificationHandler := deps.App[*handlers.NotificationHandler]()
	web.GET("/test", middleware.Auth(), notificationHandler.TestSendNotification)

}
