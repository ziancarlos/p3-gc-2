package main

import (
	"context"
	"fmt"
	"log"
	"p3-graded-challenge-2-ziancarlos/config"
	"p3-graded-challenge-2-ziancarlos/controllers"
	_ "p3-graded-challenge-2-ziancarlos/docs"
	"p3-graded-challenge-2-ziancarlos/middleware"
	"p3-graded-challenge-2-ziancarlos/repository"
	"p3-graded-challenge-2-ziancarlos/scheduler"
	"p3-graded-challenge-2-ziancarlos/service"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Shopping & Payment API
// @version 1.0
// @description This is a shopping and payment service API with gRPC and REST support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9051
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize JWT
	middleware.InitJWT(cfg.JWTSecret)

	// Connect to MongoDB
	client, err := config.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup repositories
	productCollection := config.GetCollection(client, cfg.ShoppingDBName, "products")
	paymentCollection := config.GetCollection(client, cfg.PaymentDBName, "payments")

	productRepo := repository.NewProductRepository(productCollection)
	paymentRepo := repository.NewPaymentRepository(paymentCollection)

	// Setup services
	productService := service.NewProductService(productRepo)
	paymentService := service.NewPaymentService(paymentRepo)

	// Setup controllers
	productController := controllers.NewProductController(productService)
	paymentController := controllers.NewPaymentController(paymentService)
	authController := controllers.NewAuthController()

	// Setup and start cleanup scheduler (runs every 24 hours)
	cleanupScheduler := scheduler.NewCleanupScheduler(paymentCollection, productCollection, 24*time.Hour)
	go cleanupScheduler.Start(context.Background())

	// Setup Gin router
	router := gin.Default()

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/login", authController.Login)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTMiddleware())
		{
			// Product routes
			protected.POST("/products", productController.CreateProduct)
			protected.GET("/products", productController.GetAllProducts)
			protected.GET("/products/:id", productController.GetProductByID)
			protected.PUT("/products/:id", productController.UpdateProduct)
			protected.DELETE("/products/:id", productController.DeleteProduct)

			// Payment routes
			protected.POST("/payments", paymentController.CreatePayment)
			protected.GET("/payments", paymentController.GetAllPayments)
			protected.GET("/payments/:id", paymentController.GetPaymentByID)
			protected.DELETE("/payments/:id", paymentController.DeletePayment)
		}
	}

	// Start server
	address := fmt.Sprintf(":%s", cfg.PortShopping)
	log.Printf("HTTP server listening on %s", address)
	log.Printf("Swagger documentation available at http://localhost%s/swagger/index.html", address)

	if err := router.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

