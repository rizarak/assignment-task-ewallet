package pgsql_gorm

import (
	"assignment-task-ewallet/user/entity"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type GormDBIface interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

// mendefinisikan interface untuk repository pengguna
type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

type userRepository struct {
	stDb GormDBIface
}

func NewUserRepository(db GormDBIface) IUserRepository {
	return &userRepository{stDb: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	if err := r.stDb.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("Error creating user: %v\n", err)
		return entity.User{}, err
	}

	return *user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	if err := r.stDb.WithContext(ctx).Select("id", "name", "email", "password", "created_at", "updated_at").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		log.Printf("Error getting user by id: %v\n", err)
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	var existingUser entity.User
	if err := r.stDb.WithContext(ctx).Select("id", "name", "email", "password", "created_at", "updated_at").First(&existingUser, id).Error; err != nil {
		log.Printf("Error getting user by id: %v\n", err)
		return entity.User{}, err
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Password = user.Password
	if err := r.stDb.WithContext(ctx).Save(&existingUser).Error; err != nil {
		return entity.User{}, err
	}

	return existingUser, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	if err := r.stDb.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}

	return nil
}

// mengembalikan semua pengguna
func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	if err := r.stDb.WithContext(ctx).Select("id", "name", "email", "password", "created_at", "updated_at").Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil
		}

		log.Printf("Error getting all users: %v\n", err)
		return nil, err
	}

	return users, nil
}
