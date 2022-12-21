package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user = entity.User{}
	if err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.User{}, nil
		} else {
			return entity.User{}, nil
		}
	}
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user = entity.User{}
	if err := r.db.WithContext(ctx).Table("users").Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.User{}, nil
		} else {
			return entity.User{}, nil
		}
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return entity.User{}, nil
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := r.db.WithContext(ctx).Table("users").Where("id = ?", user.ID).Update("user", user).Error; err != nil {
		return entity.User{}, nil
	}
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Scan(&entity.User{}).Delete("id=?", id).Error; err != nil {
		return err
	}
	return nil
}
