package services

import (
	"context"
	"errors"

	"github.com/berkkaradalan/AwsGo-Storage/config"
	"github.com/berkkaradalan/AwsGo-Storage/models"
	"github.com/berkkaradalan/AwsGo-Storage/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
	authconfig *config.AuthConfig
}

func NewUserService(userRepo *repositories.UserRepository, authConfig *config.AuthConfig) *UserService {
	return &UserService{
		userRepo: userRepo,
		authconfig: authConfig,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	user.UserPassword = ""

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, req.UserEmail)
	if existingUser != nil {
		return nil, errors.New("email is already in use")
	}

	existingUser, _ = s.userRepo.GetUserByUserName(ctx, req.UserName)

	if existingUser != nil {
		return nil, errors.New("username is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), 10)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		UserID: uuid.New().String(),
		UserName: req.UserName,
		UserEmail: req.UserEmail,
		UserPassword: string(hashedPassword),
	}

	user.SetTimestamps()

	createdUser, err := s.userRepo.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil { 
		return "", nil, err
	}

	if user == nil {
        return "", nil, errors.New("user not found")
    }

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token, err := s.authconfig.GenerateToken(user.UserID, user.UserName, user.UserEmail, user.CreatedAt, user.UpdatedAt)

	if err != nil { 
		return "", nil, err
	}

	return token, user, nil
}