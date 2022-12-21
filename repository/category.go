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
	var category []entity.Category
	if err := r.db.WithContext(ctx).Table("id =?", id).Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return 0, err
	}
	return category.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	if err := r.db.WithContext(ctx).Create(categories).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Category{}, nil
		}
		return entity.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	if err := r.db.WithContext(ctx).Model(category).Updates(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(entity.Category{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
