package repositories

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/models"
	"github.com/google/uuid"
)

const StorageTable = "storage"

type StorageRepository struct {
	s3Service 		*config.S3BucketService
	dynamoService	*config.DynamoDBService
	bucketName		string
}

func NewStorageRepository(dynamoservice *config.DynamoDBService, s3service *config.S3BucketService) *StorageRepository {
	return &StorageRepository{
		s3Service: s3service,
		dynamoService: dynamoservice,
		bucketName: config.LoadEnv().S3_BUCKET_NAME,
	}
}

func (r *StorageRepository) UploadFile(ctx context.Context, userID string, fileName string, fileSize int64, contentType string, fileData io.Reader, description *string) (*models.StorageObject, error) {
	objectID := uuid.New().String()
	s3Key := fmt.Sprintf("users/%s/%s", userID, objectID)

	_, err := r.s3Service.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(s3Key),
		Body:        fileData,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		return nil, err
	}

	storageObj := &models.StorageObject{
		ObjectID:    objectID,
		UserID:      userID,
		FileName:    fileName,
		FileSize:    fileSize,
		ContentType: contentType,
		S3Key:       s3Key,
		S3Bucket:    r.bucketName,
		UploadedAt:  time.Now(),
		UpdatedAt:   time.Now(),
		Description: description,
	}

	item, err := attributevalue.MarshalMap(storageObj)
	if err != nil {
		log.Printf("Failed to marshal storage object: %v", err)
		return nil, err
	}

	_, err = r.dynamoService.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(StorageTable),
		Item:      item,
	})

	if err != nil {
		log.Printf("Failed to save metadata to DynamoDB: %v", err)
		return nil, err
	}

	return storageObj, nil
}
