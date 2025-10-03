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

func NewUserRepository(service *config.DynamoDBService) *UserRepository {
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
		return nil, fmt.Errorf("User with id : %v not found", userId)
	}

	err = attributevalue.UnmarshalMap(result.Item, &user)

	if err != nil {
		log.Printf("User unmarshall failed here's why : %v", err)
		return nil, err
	}

	return &user, err

}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	log.Printf("User before marshal: %+v", user)
	// user.SetTimestamps()

	item, err := attributevalue.MarshalMap(*user)

	if err != nil {
		return nil, err
	}

	_, err = r.service.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(UsersTable),
		Item:      item,
	})

	log.Printf("Attempting to insert into table: %s", UsersTable)
	log.Printf("Marshaled item: %+v", item)

	if err != nil {
		return nil, err
	}

	return user, nil
}
