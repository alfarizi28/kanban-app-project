package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	var res []entity.Category
	tx := r.db.Raw("SELECT * FROM categories WHERE user_id = ?", id).Scan(&res)
	if tx.Error != nil {
		return []entity.Category{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return []entity.Category{}, nil
	}
	return res, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	tx := r.db.Create(category)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return category.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	tx := r.db.Create(&categories)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var res entity.Category
	tx := r.db.Raw("SELECT * FROM categories WHERE id = ?", id).Scan(&res)
	if tx.Error != nil {
		return entity.Category{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.Category{}, nil
	}
	return res, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	tx := r.db.Model(&entity.Category{}).Updates(category)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	tx := r.db.Where("id = ?", id).Delete(&entity.Category{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
