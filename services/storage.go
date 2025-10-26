package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/models"
	"github.com/berkkaradalan/AwsGo-Storage/repositories"
)

type StorageService struct {
	storageRepo *repositories.StorageRepository
	authconfig *config.AuthConfig
}

func NewStorageService(storageRepo *repositories.StorageRepository, authConfig *config.AuthConfig) *StorageService{
	return &StorageService{
		storageRepo: storageRepo,
		authconfig: authConfig,
	}
}

func (s *StorageService) UploadFile(ctx context.Context, userID string, file *multipart.FileHeader, description *string) (*models.UploadFileResponse, error) {
	if file == nil {
		return nil, fmt.Errorf("file is required.")
	}

	if file.Size > 50*1024*1024 { // 50 MB limit
		return nil, fmt.Errorf("file size exceeds 50 MB limit")
	}

	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"application/pdf": true,
	}

	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("file type not allowed. Allowed: JPEG, PNG, GIF, WebP, PDF")
	}

	src, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return nil, fmt.Errorf("failed to open file")
	}
	defer src.Close()

	storageObj, err := s.storageRepo.UploadFile(ctx, userID, file.Filename, file.Size, contentType, src, description)
	if err != nil {
		return nil, err
	}

	response := &models.UploadFileResponse{
		ObjectID:    storageObj.ObjectID,
		FileName:    storageObj.FileName,
		FileSize:    storageObj.FileSize,
		ContentType: storageObj.ContentType,
		UploadedAt:  storageObj.UploadedAt,
		Description: storageObj.Description,
		Message:     "File uploaded successfully",
	}

	return response, nil
}

func (s *StorageService) ListFiles(ctx context.Context, userID string) (*models.ListStorageObjectsResponse, error) {
	if userID == "" {
		return nil, errors.New("file ID cannot be empty")
	}

	files, err := s.storageRepo.ListFiles(ctx, userID)

	if err != nil {
		return nil, err
	}

	for i := range files.Data {
		previewURL, err := s.storageRepo.GeneratePresignedURL(
			ctx, 
			files.Data[i].S3Key,
			files.Data[i].ContentType,
			30*time.Minute,
		)
		if err != nil {
			log.Printf("Failed to generate preview URL for %s: %v", files.Data[i].ObjectID, err)
			continue
		}
		files.Data[i].PreviewURL = previewURL
	}

	return files, nil
}

func (s *StorageService) DownloadFile(ctx context.Context, fileID string, userID string) ([]byte, error) {
	if fileID == "" {
		return nil, errors.New("file ID cannot be empty")
	}

	fileData, err := s.storageRepo.DownloadFile(ctx, fileID, userID)

	if err != nil {
        return nil, err
    }

    if fileData == nil {
        return nil, errors.New("file not found")
    }

    return fileData, nil
}