package service

import (
	"assignment-task-ewallet/user/entity"
	"context"
	"fmt"

	pgsqlGorm "assignment-task-ewallet/user/repository/pgsql_gorm"
)

// IUserService mendefinisikan interface untuk layanan pengguna
type IUserService interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

// userService adalah implementasi dari IUserService yang menggunakan IUserRepository
type userService struct {
	userRepo pgsqlGorm.IUserRepository
}

// NewUserService membuat instance baru dari userService
func NewUserService(userRepo pgsqlGorm.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

// CreateUser membuat pengguna baru
func (s *userService) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("Error creating user: %v", err)
	}

	return createdUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	selectedUser, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("User not found: %v", err)
	}

	return selectedUser, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	updatedUser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("Error updating user: %v", err)
	}

	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}

	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting all users: %v", err)
	}

	return users, nil
}
