package models

import "time"

type StorageObject struct {
    ObjectID      string    `dynamodbav:"ObjectID"`
    UserID        string    `dynamodbav:"UserID"`
    FileName      string    `dynamodbav:"FileName"`
    FileSize      int64     `dynamodbav:"FileSize"`
    ContentType   string    `dynamodbav:"ContentType"`
    S3Key         string    `dynamodbav:"S3Key"`
    S3Bucket      string    `dynamodbav:"S3Bucket"`
    UploadedAt    time.Time `dynamodbav:"UploadedAt"`
    UpdatedAt     time.Time `dynamodbav:"UpdatedAt"`
    Description   *string    `dynamodbav:"Description"`
}

type UploadFileRequest struct {
    Description *string `form:"description"`
}

type UploadFileResponse struct {
    ObjectID    string    `json:"objectId"`
    FileName    string    `json:"fileName"`
    FileSize    int64     `json:"fileSize"`
    ContentType string    `json:"contentType"`
    UploadedAt  time.Time `json:"uploadedAt"`
    Description *string   `json:"description,omitempty"`
    Message     string    `json:"message"`
}

type ListStorageObjectsResponse struct {
    Success  bool             `json:"success"`
    Message  string           `json:"message"`
    Data     []StorageObject  `json:"data"`
    Count    int              `json:"count"`
    // NextToken *string         `json:"nextToken,omitempty"`
}