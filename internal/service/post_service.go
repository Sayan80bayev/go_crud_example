package service

import (
	"errors"
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

func (ps *PostService) UpdatePost(content string, title string, userId uint, postId uint) error {
	post, err2 := validateUpdate(content, title, userId, postId, ps)
	if err2 != nil {
		return err2
	}

	post.Content = content
	post.Title = title

	return ps.postRepo.UpdatePost(post)
}

func (ps *PostService) DeletePost(postId uint, userId uint) error {
	post, err := ps.postRepo.GetPostByID(postId)
	if err != nil {
		return errors.New("post not found")
	}
	if post.User.ID != userId {
		return errors.New("user not allowed")
	}

	return ps.postRepo.DeletePost(postId)
}

func (uc *PostService) GetPostByID(id uint) (*models.Post, error) {
	return uc.postRepo.GetPostByID(id)
}

func validateUpdate(content string, title string, userId uint, postId uint, ps *PostService) (*models.Post, error) {
	post, err := ps.postRepo.GetPostByID(postId)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.User.ID != userId {
		return nil, errors.New("user not allowed")
	}

	isTitleEqual := title == post.Title
	isContentEqual := content == post.Content
	if isTitleEqual && isContentEqual {
		return nil, errors.New("nothing to update")
	}
	return post, nil
}
