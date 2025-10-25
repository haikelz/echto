package service

import (
	"echto/internal/entity"
	"echto/internal/model"
	"echto/internal/repository"
	"echto/pkg/logger"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req *model.UserCreateRequest) (*model.UserResponse, error)
	GetUser(id uint) (*model.UserResponse, error)
	GetUsers(page, limit int) (*model.UserListResponse, error)
	UpdateUser(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(req *model.UserCreateRequest) (*model.UserResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		return nil, errors.New("failed to process password")
	}

	// Create user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create user")
		return nil, errors.New("failed to create user")
	}

	// Return response
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUser(id uint) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errors.New("user not found")
		}
		logger.Log.Error().Err(err).Msg("Failed to get user")
		return nil, errors.New("failed to get user")
	}

	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) GetUsers(page, limit int) (*model.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := s.userRepo.GetAll(page, limit)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get users")
		return nil, errors.New("failed to get users")
	}

	// Convert to response format
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = model.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return &model.UserListResponse{
		Users: userResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *userService) UpdateUser(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errors.New("user not found")
		}
		logger.Log.Error().Err(err).Msg("Failed to get user")
		return nil, errors.New("failed to get user")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email already exists (excluding current user)
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to update user")
		return nil, errors.New("failed to update user")
	}

	// Return response
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New("user not found")
		}
		logger.Log.Error().Err(err).Msg("Failed to get user")
		return errors.New("failed to get user")
	}

	// Delete user
	if err := s.userRepo.Delete(id); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete user")
		return errors.New("failed to delete user")
	}

	return nil
}
