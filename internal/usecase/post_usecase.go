package usecase

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
)

type PostUsecase struct {
	repo *repository.PostRepository
}

func NewPostUsecase(repo *repository.PostRepository) *PostUsecase {
	return &PostUsecase{repo: repo}
}

func (uc *PostUsecase) CreatePost(post *models.Post) error {
	return uc.repo.Create(post)
}

func (uc *PostUsecase) GetPosts() ([]models.Post, error) {
	return uc.repo.GetAll()
}

func (uc *PostUsecase) GetPostByID(id uint) (*models.Post, error) {
	return uc.repo.GetByID(id)
}

func (uc *PostUsecase) UpdatePost(post *models.Post) error {
	return uc.repo.Update(post)
}

func (uc *PostUsecase) DeletePost(id uint) error {
	return uc.repo.Delete(id)
}
