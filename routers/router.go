package routers

import (
	"github.com/berkkaradalan/AwsGo-Storage/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine{
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
		routes.POST("/user", userHandler.CreateUser)
	}

	return router
}