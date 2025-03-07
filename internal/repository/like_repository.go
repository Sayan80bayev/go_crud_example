package repository

import (
	"go_crud_example/internal/models"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db}
}

func (r *LikeRepository) CreateLike(like models.Like) error {
	return r.db.Create(&like).Error
}

func (r *LikeRepository) DeleteLike(id uint) error {
	return r.db.Delete(&models.Like{}, id).Error
}
