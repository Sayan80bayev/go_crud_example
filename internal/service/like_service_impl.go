package service

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
)

type likeServiceImpl struct {
	repository *repository.LikeRepository
}

func NewLikeService(repository *repository.LikeRepository) LikeService {
	return &likeServiceImpl{repository}
}

func (l likeServiceImpl) LikePost(userID, postID uint) error {
	like := models.Like{
		UserID: userID,
		PostID: postID,
	}

	return l.repository.CreateLike(like)
}

func (l likeServiceImpl) UnlikePost(likeID uint) error {
	return l.repository.DeleteLike(likeID)
}
