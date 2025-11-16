package api

import (
	"github.com/dwfennell/monorepo-scaffold/backend/internal/database"
	"github.com/dwfennell/monorepo-scaffold/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *database.DB) {
	userRepo := repository.NewUserRepository(db)
	authHandler := NewAuthHandler(userRepo)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/me", authHandler.GetCurrentUser)
		}
	}
}
