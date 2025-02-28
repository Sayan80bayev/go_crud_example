package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go_crud_example/internal/config"
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/response"
	"go_crud_example/pkg/s3"
)

type PostService struct {
	postRepo *repository.PostRepository
	minio    *minio.Client
}

func NewPostService(postRepo *repository.PostRepository, minioClient *minio.Client) *PostService {
	return &PostService{postRepo, minioClient}
}

func (ps *PostService) CreatePost(c *gin.Context, title, content string, userID uint, categoryID uint, cfg *config.Config) error {
	// Upload the image to MinIO
	imageURL, err := s3.UploadFile(c, cfg, ps.minio)
	if err != nil {
		return err
	}

	// Create post with image URL
	post := &models.Post{
		Title:      title,
		Content:    content,
		UserID:     userID,
		CategoryID: categoryID,
		ImageURL:   imageURL, // Store image URL in DB
	}

	return ps.postRepo.CreatePost(post)
}

func (ps *PostService) GetPosts() ([]response.PostResponse, error) {
	posts, err := ps.postRepo.GetPosts()
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
			Category: response.CategoryResponse{
				Id:   post.Category.ID,
				Name: post.Category.Name,
			},
			ImageURL: &post.ImageURL,
		})
	}
	return postResponses, nil
}

func (ps *PostService) UpdatePost(c *gin.Context, content string, title string, userId uint, postId uint, cfg *config.Config) error {
	post, err2 := validateUpdate(userId, postId, ps)
	if err2 != nil {
		return err2
	}

	imageURL, err := s3.UploadFile(c, cfg, ps.minio)
	if err != nil {
		return err
	}

	post.Content = content
	post.Title = title
	post.ImageURL = imageURL

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

func (ps *PostService) GetPostByID(id uint) (*models.Post, error) {
	return ps.postRepo.GetPostByID(id)
}

func validateUpdate(userId uint, postId uint, ps *PostService) (*models.Post, error) {
	post, err := ps.postRepo.GetPostByID(postId)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.User.ID != userId {
		return nil, errors.New("user not allowed")
	}

	return post, nil
}
