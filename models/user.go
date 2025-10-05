package models

import "time"

type User struct {
	UserID       string `json:"user_id" dynamodbav:"UserID"`
	UserName     string `json:"user_name" dynamodbav:"UserName"`
	UserEmail    string `json:"user_email" dynamodbav:"UserEmail"`
	UserPassword string `json:"user_password,omitempty" dynamodbav:"UserPassword"`
	CreatedAt    int64  `json:"created_at" dynamodbav:"CreatedAt"`
	UpdatedAt    int64  `json:"updated_at" dynamodbav:"UpdatedAt"`
}

type UserResponse struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type LoginRequest struct {
	UserEmail	 string `json:"user_email" binding:"required,email"`
	UserPassword string `json:"user_password" binding:"required,min=8,max=50"`
}

type CreateUserRequest struct {
	UserName     string `json:"user_name" binding:"required,min=3,max=50"`
	UserEmail    string `json:"user_email" binding:"required,email"`
	UserPassword string `json:"user_password" binding:"required,min=8,max=50"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		UserID:    u.UserID,
		UserName:  u.UserName,
		UserEmail: u.UserEmail,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) SetTimestamps() {
	now := time.Now().Unix()
	if u.CreatedAt == 0 {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}
