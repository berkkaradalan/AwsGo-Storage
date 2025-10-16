package routers

import (
	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/handlers"
	"github.com/berkkaradalan/AwsGo-Storage/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, storageHandler *handlers.StorageHandler, env config.Env, authConfig *config.AuthConfig) *gin.Engine{
	router := gin.Default()
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"message": "Server is running",
		})
	})

	routes := router.Group("/api/v1") 
	{
		routes.GET("/user/:id", userHandler.GetUserByID)
		routes.POST("/user/register", userHandler.CreateUser)
		routes.POST("/user/login", userHandler.Login)
		routes.POST("/storage/upload", storageHandler.UploadFile)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authConfig))
	{
		protected.GET("/user/me", userHandler.GetProfile)
	}

	return router
}