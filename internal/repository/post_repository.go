package repository

import (
	"go_crud_example/internal/models"
	"gorm.io/gorm"
)

// PostRepository — интерфейс работы с постами
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository — конструктор репозитория
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create добавляет новый пост
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// GetAll возвращает все посты
func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Find(&posts).Error
	return posts, err
}

// GetByID возвращает пост по ID
func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update обновляет пост
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete удаляет пост
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
