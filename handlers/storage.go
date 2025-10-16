package handlers

import (
	"net/http"

	"github.com/berkkaradalan/AwsGo-Storage/middleware"
	"github.com/berkkaradalan/AwsGo-Storage/services"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	storageService *services.StorageService
}

func NewStorageHandler(storageService *services.StorageService) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
	}
}

func (h *StorageHandler) UploadFile(c *gin.Context) {
	userData := middleware.GetCurrentClaims(c)
	if userData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userData.UserID

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	description := c.PostForm("description")
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	response, err := h.storageService.UploadFile(c.Request.Context(), userID, file, descPtr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}