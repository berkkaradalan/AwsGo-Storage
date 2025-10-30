package config

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type DynamoDBService struct {
	Client *dynamodb.Client
}

type S3BucketService struct {
	Client *s3.Client
}

func (client *DynamoDBService) TableExists(ctx context.Context, tableName string) (bool, error) {
	exists := true
	_, err := client.Client.DescribeTable(
		ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
	)
	if err != nil {
		var notFoundEx *dynamotypes.ResourceNotFoundException
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

func (client *DynamoDBService) CreateTable(ctx context.Context, createTableInput dynamodb.CreateTableInput, tableName string) (*dynamotypes.TableDescription, error) {
	var tableDesc *dynamotypes.TableDescription
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
		KeySchema: []dynamotypes.KeySchemaElement{
			{
				AttributeName: aws.String("UserID"),
				KeyType:       dynamotypes.KeyTypeHash, // Partition key
			},
		},
		AttributeDefinitions: []dynamotypes.AttributeDefinition{
			{
				AttributeName: aws.String("UserID"),
				AttributeType: dynamotypes.ScalarAttributeTypeS, // String
			},
			{
				AttributeName: aws.String("UserEmail"),
				AttributeType: dynamotypes.ScalarAttributeTypeS, // String
			},
		},
		BillingMode: dynamotypes.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []dynamotypes.GlobalSecondaryIndex{
			{
				IndexName: aws.String("UserEmailIndex"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{
						AttributeName: aws.String("UserEmail"),
						KeyType:       dynamotypes.KeyTypeHash,
					},
				},
				Projection: &dynamotypes.Projection{
					ProjectionType: dynamotypes.ProjectionTypeAll,
				},
			},
		},
	}
}

func CreateStorageTableInput() dynamodb.CreateTableInput {
	return dynamodb.CreateTableInput{
		TableName: aws.String("storage"),
		KeySchema: []dynamotypes.KeySchemaElement{
			{
				AttributeName: aws.String("ObjectID"),
				KeyType:       dynamotypes.KeyTypeHash, // Partition key
			},
		},
		AttributeDefinitions: []dynamotypes.AttributeDefinition{
			{
				AttributeName: aws.String("ObjectID"),
				AttributeType: dynamotypes.ScalarAttributeTypeS, // String
			},
			{
				AttributeName: aws.String("UserID"),
				AttributeType: dynamotypes.ScalarAttributeTypeS, // String
			},
			{
				AttributeName: aws.String("UploadedAt"),
				AttributeType: dynamotypes.ScalarAttributeTypeS, // String (ISO 8601 format)
			},
		},
		BillingMode: dynamotypes.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []dynamotypes.GlobalSecondaryIndex{
			{
				IndexName: aws.String("UserIDIndex"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{
						AttributeName: aws.String("UserID"),
						KeyType:       dynamotypes.KeyTypeHash,
					},
				},
				Projection: &dynamotypes.Projection{
					ProjectionType: dynamotypes.ProjectionTypeAll,
				},
			},
			{
				IndexName: aws.String("UserIDUploadedAtIndex"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{
						AttributeName: aws.String("UserID"),
						KeyType:       dynamotypes.KeyTypeHash,
					},
					{
						AttributeName: aws.String("UploadedAt"),
						KeyType:       dynamotypes.KeyTypeRange, // Sort key
					},
				},
				Projection: &dynamotypes.Projection{
					ProjectionType: dynamotypes.ProjectionTypeAll,
				},
			},
		},
	}
}

func ConnectDatabase() *DynamoDBService {
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

	storageTableCheck, err := service.TableExists(context.TODO(), "storage")

	if err != nil {
		panic(err)
	}

	if !storageTableCheck {
		log.Println("Creating storage table")
		storageTableInput := CreateStorageTableInput()
		_, err := service.CreateTable(context.Background(), storageTableInput, "storage")
		if err != nil {
			panic(err)
		}
	}

	return service
}

func (client *S3BucketService) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := client.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true

	log.Printf("headbucket after")

	if err != nil {
		var apiError smithy.APIError

		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *s3types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		} else {
			log.Printf("Couldn't determine existence of bucket %v. Here's why: %v\n", bucketName, err)
		}
		exists = false
	}

	return exists, err
}

func (client *S3BucketService) CreateBucket(ctx context.Context, bucketName string, region string) error {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	if region != "us-east-1" {
		input.CreateBucketConfiguration = &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(region),
		}
	}

	_, err := client.Client.CreateBucket(ctx, input)

	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exist *types.BucketAlreadyExists

		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s. \n", bucketName)
			err = owned
		} else if errors.As(err, &exist) {
			log.Printf("Bucket %s already exist.\n", bucketName)
		}
	} else {
		err = s3.NewBucketExistsWaiter(client.Client).Wait(
			ctx, &s3.HeadBucketInput{
				Bucket: aws.String(bucketName),
			}, time.Minute)

		if err != nil {
			log.Printf("Failed to attempt to wait for bucket %s to exist. \n", bucketName)
			panic(err)
		}
	}
	return err
}

func ConnectS3Bucket(env *Env) *S3BucketService {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Println(err)
		panic(err)
	}

	client := s3.NewFromConfig(cfg)

	service := &S3BucketService{
		Client: client,
	}

	userStorageBucket, err := service.BucketExists(context.Background(), env.S3_BUCKET_NAME)

	if err != nil {
		panic(err)
	}

	if !userStorageBucket {
		log.Printf("Creating user bucket with name : %v ", env.S3_BUCKET_NAME)
		log.Println(cfg.Region)
		err = service.CreateBucket(context.Background(), env.S3_BUCKET_NAME, cfg.Region)

		if err != nil {
			panic(err)
		}
	}

	return service
}
