package service

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
)

type LikeService struct {
	repository *repository.LikeRepository
}

func NewLikeService(repository *repository.LikeRepository) *LikeService {
	return &LikeService{repository}
}

func (l LikeService) LikePost(userID, postID uint) error {
	like := models.Like{
		UserID: userID,
		PostID: postID,
	}

	return l.repository.CreateLike(like)
}

func (l LikeService) UnlikePost(likeID uint) error {
	return l.repository.DeleteLike(likeID)
}
