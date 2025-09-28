package config

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBService struct {
	Client *dynamodb.Client
}

func (client *DynamoDBService) TableExists(ctx context.Context, tableName string) (bool, error) {
	exists := true
	_, err := client.Client.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			log.Printf("Table %v does not exist.\n", tableName)
			err = nil
		} else {
			log.Printf("Couldn't determine existence of table %v. Here's why: %v\n", tableName, err)
		}
		exists = false
	}
	return exists, err
}

func (client *DynamoDBService) CreateTable(ctx context.Context, createTableInput dynamodb.CreateTableInput, tableName string) (*types.TableDescription, error){
	var tableDesc *types.TableDescription
	table, err := client.Client.CreateTable(ctx, &createTableInput)

	if err != nil {
		log.Printf("Couldn't create table. Here's why: %v\n", err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(client.Client)
		err = waiter.Wait(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
		log.Printf("Ccreating table test")
	}
	return tableDesc, err
}

func CreateUserTableInput() dynamodb.CreateTableInput {
	return dynamodb.CreateTableInput{
		TableName: aws.String("user"),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("UserID"),
				KeyType:       types.KeyTypeHash, // Partition key
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("UserID"),
				AttributeType: types.ScalarAttributeTypeS, // String
			},
			{
				AttributeName: aws.String("UserEmail"),
				AttributeType: types.ScalarAttributeTypeS, // String
			},
		},
		BillingMode: types.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("UserEmailIndex"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("UserEmail"),
						KeyType:       types.KeyTypeHash,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
	}
}

func ConnectDatabase() *DynamoDBService{
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Printf("Failed to connect to database %v", err)
		panic(err)
	}

	client := dynamodb.NewFromConfig(cfg)
	service := &DynamoDBService{Client: client}


	userTableCheck, err := service.TableExists(context.TODO(), "user")

	if err != nil {
		panic(err)
	}

	if !userTableCheck {
		log.Println("Creating user table")
		userTableInput := CreateUserTableInput()
		_, err := service.CreateTable(context.Background(), userTableInput, "user")
		if err != nil {
			panic(err)
		}
	}

	return service
}
