package usecase

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
)

type PostUsecase struct {
	postRepo *repository.PostRepository
}

func NewPostUsecase(postRepo *repository.PostRepository) *PostUsecase {
	return &PostUsecase{postRepo}
}

func (uc *PostUsecase) CreatePost(title, content string, userID uint) error {
	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  userID, // Указываем автора
	}
	return uc.postRepo.CreatePost(post)
}

func (uc *PostUsecase) GetPosts() ([]models.Post, error) {
	return uc.postRepo.GetPosts()
}

func (uc *PostUsecase) GetPostByID(id uint) (*models.Post, error) {
	return uc.postRepo.GetPostByID(id)
}
