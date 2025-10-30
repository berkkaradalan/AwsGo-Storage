package handlers

import (
	"fmt"
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

func (h *StorageHandler) ListFiles(c *gin.Context) {
	userData := middleware.GetCurrentClaims(c)
	userID := userData.UserID

	files, err := h.storageService.ListFiles(c, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *StorageHandler) DownloadFile(c *gin.Context) {
	userData := middleware.GetCurrentClaims(c)
	userID := userData.UserID
	fileID := c.Param("id")
	
	fileData, err := h.storageService.DownloadFile(c.Request.Context(), fileID, userID)

	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"error" : "Error while downloading file."})
		return 
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileID))
    c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Length", fmt.Sprintf("%d", len(fileData)))
    
    c.Data(http.StatusOK, "application/octet-stream", fileData)

}

func (h *StorageHandler) DeleteFile(c *gin.Context) {
	userData := middleware.GetCurrentClaims(c)
	userID := userData.UserID
	fileID := c.Param("id")

	deleteMessage, err := h.storageService.DeleteFile(c, userID, fileID)

	if err != nil ||  deleteMessage == nil{
		c.JSON(http.StatusExpectationFailed, gin.H{"error": "Error while deleting file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message":deleteMessage})
}

func (h *StorageHandler) GetDashboardMetrics(c *gin.Context) {
	userData := middleware.GetCurrentClaims(c)
	userID := userData.UserID

	dashboardMetrics, err := h.storageService.GetDashboardMetrics(c, userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while getting dashboard metrics"})
	}

	c.JSON(http.StatusOK, dashboardMetrics)
}