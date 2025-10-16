package repositories

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/models"
)

const UsersTable = "user"

type UserRepository struct {
	service *config.DynamoDBService
}

func NewUserRepository(service *config.DynamoDBService, s3Service *config.S3BucketService) *UserRepository {
	return &UserRepository{
		service: service,
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, userId string) (*models.User, error) {
	var user models.User
	result, err := r.service.Client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberS{
				Value: userId,
			},
		},
		TableName: aws.String(UsersTable),
	})

	if err != nil {
		log.Printf("Couldn't get user with id : %v, Here's what happened : %v", userId, err)
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("user with id : %v not found", userId)
	}

	err = attributevalue.UnmarshalMap(result.Item, &user)

	if err != nil {
		log.Printf("User unmarshall failed here's why : %v", err)
		return nil, err
	}

	return &user, err

}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	item, err := attributevalue.MarshalMap(*user)

	if err != nil {
		return nil, err
	}

	_, err = r.service.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(UsersTable),
		Item:      item,
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	result, err := r.service.Client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(UsersTable),
		IndexName: aws.String("UserEmailIndex"),
		KeyConditionExpression: aws.String("UserEmail = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: userEmail},
		},
	})

	if err != nil {
		log.Printf("couldn't get user with email: %v, error: %v", userEmail, err)
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(result.Items[0], &user); err != nil {
		log.Printf("User unmarshal failed: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUserName(ctx context.Context, userName string) (*models.User, error) {
	result, err := r.service.Client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(UsersTable),
		IndexName: aws.String("UserNameIndex"),
		KeyConditionExpression: aws.String("UserName = :username"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: userName},
		},
	})

	if err != nil {
		log.Printf("couldn't get user with username: %v, error: %v", userName, err)
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var user models.User
	if err := attributevalue.UnmarshalMap(result.Items[0], &user); err != nil {
		log.Printf("User unmarshal failed: %v", err)
		return nil, err
	}

	return &user, nil
}