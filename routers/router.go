package routers

import (
	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/handlers"
	"github.com/berkkaradalan/AwsGo-Storage/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler, storageHandler *handlers.StorageHandler, env config.Env, authConfig *config.AuthConfig) *gin.Engine{
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	
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
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authConfig))
	{
		protected.GET("/user/me", userHandler.GetProfile)
		protected.POST("/storage/upload", storageHandler.UploadFile)
		protected.GET("/storage/files", storageHandler.ListFiles)
		protected.GET("/storage/files/:id/download", storageHandler.DownloadFile)
	}

	return router
}