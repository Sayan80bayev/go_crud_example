package usecase

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/response"
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

func (uc *PostUsecase) GetPosts() ([]response.PostResponse, error) {
	posts, err := uc.postRepo.GetPosts()
	if err != nil {
		return nil, err
	}

	var postResponses []response.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, response.PostResponse{
			ID:    post.ID,
			Title: post.Title,
			Author: response.UserResponse{
				ID:       post.User.ID,
				Username: post.User.Username,
			},
		})
	}
	return postResponses, nil
}

func (uc *PostUsecase) GetPostByID(id uint) (*models.Post, error) {
	return uc.postRepo.GetPostByID(id)
}
