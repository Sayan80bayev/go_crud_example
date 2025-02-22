package service

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/response"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo}
}

func (uc *PostService) CreatePost(title, content string, userID uint) error {
	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  userID, // Указываем автора
	}
	return uc.postRepo.CreatePost(post)
}

func (uc *PostService) GetPosts() ([]response.PostResponse, error) {
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

func (uc *PostService) GetPostByID(id uint) (*models.Post, error) {
	return uc.postRepo.GetPostByID(id)
}
