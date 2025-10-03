package handlers

import (
	"net/http"

	"github.com/berkkaradalan/AwsGo-Storage/models"
	"github.com/berkkaradalan/AwsGo-Storage/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(c, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": user.ToResponse(),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": user.ToResponse(),
	})
}