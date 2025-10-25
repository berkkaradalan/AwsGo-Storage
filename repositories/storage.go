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