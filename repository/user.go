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
	var res entity.User
	tx := r.db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&res)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.User{}, nil
	}
	return res, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var res entity.User
	tx := r.db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&res)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.User{}, nil
	}
	return res, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	tx := r.db.Create(&user)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	tx := r.db.Model(&[]entity.User{}).Updates(&user)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	tx := r.db.Where("id = ?", id).Delete(&entity.User{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
