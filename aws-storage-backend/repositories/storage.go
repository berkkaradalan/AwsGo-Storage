package repositories

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (r *StorageRepository) ListFiles(ctx context.Context, userID string) (*models.ListStorageObjectsResponse, error) {
	result, err := r.dynamoService.Client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(StorageTable),
		IndexName: aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID},
		},
	})

	if err != nil {
		log.Printf("couldn't get user files with userID : %v, error : %v", userID, err)
		return nil, err
	}

	if len(result.Items) == 0 {
		return &models.ListStorageObjectsResponse{
			Success: true,
			Message: "No files found",
			Data:    []models.StorageObject{},
			Count:   0,
		}, nil
	}

	var files []models.StorageObject

	err = attributevalue.UnmarshalListOfMaps(result.Items, &files)

	if err != nil {
		log.Printf("failed to unmarshal dynamodb items: %v", err)
		return nil, err
	}

	// var nextToken *string
	// if val, ok := result.LastEvaluatedKey["userID"]; ok{
	// 	if attr, ok := val.(*types.AttributeValueMemberS); ok {
	// 		nextToken = aws.String(attr.Value)
	// 	}
	// }

	response := &models.ListStorageObjectsResponse{
		Success: true,
		Message: "Files fetched successfully",
		Data: files,
		Count: len(files),
		// NextToken: nextToken,
	}

	return response, nil
}

func (r *StorageRepository) DownloadFile(ctx context.Context, fileID string, userID string) ([]byte, error){
	result, err := r.dynamoService.Client.GetItem(ctx, &dynamodb.GetItemInput{
        TableName: aws.String(StorageTable),
        Key: map[string]types.AttributeValue{
            "ObjectID": &types.AttributeValueMemberS{Value: fileID},
        },
    })

	if err != nil {
        log.Printf("GetItem error for fileID %s: %v", fileID, err)
        return nil, err
    }

    if result.Item == nil {
        return nil, errors.New("file not found")
    }

    var storageObj models.StorageObject
    err = attributevalue.UnmarshalMap(result.Item, &storageObj)
    if err != nil {
        log.Printf("unmarshal error: %v", err)
        return nil, err
    }

    if storageObj.UserID != userID {
        return nil, errors.New("unauthorized: file does not belong to user")
    }

    output, err := r.s3Service.Client.GetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(storageObj.S3Bucket),
        Key:    aws.String(storageObj.S3Key),
    })
    if err != nil {
        log.Printf("S3 download error: %v", err)
        return nil, err
    }
    defer output.Body.Close()

    data, err := io.ReadAll(output.Body)
    if err != nil {
        log.Printf("read error: %v", err)
        return nil, err
    }

    return data, nil
}

func (r *StorageRepository) GeneratePresignedURL(ctx context.Context, s3Key string, contentType string, expiresIn time.Duration) (string, error) {
    presignClient := s3.NewPresignClient(r.s3Service.Client)

	input := &s3.GetObjectInput{
        Bucket: aws.String(r.bucketName),
        Key:    aws.String(s3Key),
    }

	if contentType == "application/pdf" {
        input.ResponseContentType = aws.String("application/pdf")
        input.ResponseContentDisposition = aws.String("inline")
    }
    
	request, err := presignClient.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
        opts.Expires = expiresIn
    })
    
    if err != nil {
        log.Printf("Failed to generate presigned URL: %v", err)
        return "", err
    }
    
    return request.URL, nil
}

func (r *StorageRepository) DeleteFile(ctx context.Context, userID string, fileID string) (*string, error){
	result, err := r.dynamoService.Client.GetItem(ctx, &dynamodb.GetItemInput{
        TableName: aws.String(StorageTable),
        Key: map[string]types.AttributeValue{
            "ObjectID": &types.AttributeValueMemberS{Value: fileID},
        },
    })

	if err != nil {
        log.Printf("GetItem error for fileID %s: %v", fileID, err)
        return nil, err
    }

    if result.Item == nil {
        return nil, errors.New("file not found")
    }

	var storageObj models.StorageObject
    err = attributevalue.UnmarshalMap(result.Item, &storageObj)
    if err != nil {
        log.Printf("unmarshal error: %v", err)
        return nil, err
    }

	if storageObj.UserID != userID {
        return nil, errors.New("unauthorized: file does not belong to user")
    }

	_, err = r.dynamoService.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(StorageTable),
		Key: map[string]types.AttributeValue{
            "ObjectID": &types.AttributeValueMemberS{Value: fileID},
        },
	})

	if err != nil {
		log.Printf("DeleteItem error for fileID : %s reason :%v", err)
		return nil, err
	}

	_, err = r.s3Service.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key: aws.String(storageObj.S3Key),
	})

	if err != nil {
		log.Printf("DeleteObject error for fileID %s: %v", fileID, err)
		return nil, err
	}

	deletionMessage := fmt.Sprintf("Object with id :%v deleted.", fileID)

	return &deletionMessage, nil
}

func (r *StorageRepository) GetDashboardMetrics(ctx context.Context, userID string) (*models.DashboardResponse, error) {
	result, err := r.dynamoService.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(StorageTable),
		IndexName:              aws.String("UserIDIndex"),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{Value: userID},
		},
	})

	if err != nil {
		log.Printf("Failed to query user files: %v", err)
		return nil, err
	}

	var files []models.StorageObject
	if len(result.Items) > 0 {
		err = attributevalue.UnmarshalListOfMaps(result.Items, &files)
		if err != nil {
			log.Printf("Failed to unmarshal files: %v", err)
			return nil, err
		}
	}

	monthlyData := groupFilesByMonth(files)

	totalSize := int64(0)
	totalFiles := 0
	for _, usage := range monthlyData {
		totalSize += usage.TotalSize
		totalFiles += usage.FileCount
	}

	return &models.DashboardResponse{
		Success: true,
		Message: "Dashboard data fetched successfully",
		Data: models.DashboardData{
			Months: monthlyData,
			Summary: map[string]interface{}{
				"totalSizeInBytes": totalSize,
				"totalSizeInMB":    float64(totalSize) / (1024 * 1024),
				"totalSizeInGB":    float64(totalSize) / (1024 * 1024 * 1024),
				"totalFiles":       totalFiles,
			},
		},
	}, nil
}

func groupFilesByMonth(files []models.StorageObject) []models.MonthlyUsage {
	now := time.Now()

	months := make([]string, 12)
	monthNames := []string{
		"Janua", "Febru", "March", "April", "May", "June",
		"July", "Augus", "Septe", "Octob", "Novem", "Decem",
	}

	monthMap := make(map[string]*models.MonthlyUsage)

	for i := 11; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		yearMonth := date.Format("2006-01")
		months[11-i] = yearMonth

		monthIndex := int(date.Month()) - 1
		monthMap[yearMonth] = &models.MonthlyUsage{
			Month:     yearMonth,
			MonthName: monthNames[monthIndex],
			TotalSize: 0,
			FileCount: 0,
			SizeInMB:  0,
		}
	}

	for _, file := range files {
		yearMonth := file.UploadedAt.Format("2006-01")

		if usage, exists := monthMap[yearMonth]; exists {
			usage.TotalSize += file.FileSize
			usage.FileCount++
			usage.SizeInMB = float64(usage.TotalSize) / (1024 * 1024)
		}
	}

	result := make([]models.MonthlyUsage, 0, 12)
	for _, month := range months {
		if usage, exists := monthMap[month]; exists {
			result = append(result, *usage)
		}
	}
	return result
}
