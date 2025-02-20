package repository

import (
	"go_crud_example/internal/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db}
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").Find(&posts).Error
	return posts, err
}

func (r *PostRepository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	return &post, err
}

func (r *PostRepository) UpdatePost(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) DeletePost(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
