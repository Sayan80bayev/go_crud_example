package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go_crud_example/internal/mappers"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go_crud_example/internal/config"
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/response"
	"go_crud_example/pkg/s3"
)

type PostService struct {
	postRepo *repository.PostRepository
	minio    *minio.Client
	redis    *redis.Client
	producer *kafka.Producer
}

func NewPostService(postRepo *repository.PostRepository, minioClient *minio.Client, redis *redis.Client, producer *kafka.Producer) *PostService {
	return &PostService{postRepo, minioClient, redis, producer}
}

func (ps *PostService) CreatePost(c *gin.Context, title, content string, userID uint, categoryID uint, cfg *config.Config) error {
	imageURL, err := s3.UploadFile(c, cfg, ps.minio)
	if err != nil {
		return err
	}

	post := &models.Post{
		Title:      title,
		Content:    content,
		UserID:     userID,
		CategoryID: categoryID,
		ImageURL:   imageURL,
		LikeCount:  0,
	}

	err = ps.postRepo.CreatePost(post)
	if err != nil {
		return err
	}

	ps.redis.Publish(context.Background(), "posts:updates", "update")

	return nil
}

func (ps *PostService) GetPosts() ([]response.PostResponse, error) {
	ctx := context.Background()
	cachedPosts, err := ps.redis.Get(ctx, "posts:list").Result()
	if err == nil {
		var postResponses []response.PostResponse
		json.Unmarshal([]byte(cachedPosts), &postResponses)
		return postResponses, nil
	}

	posts, err := ps.postRepo.GetPosts()
	if err != nil {
		return nil, err
	}

	postResponses := mappers.MapPostsToResponse(posts)
	jsonData, _ := json.Marshal(postResponses)
	ps.redis.Set(ctx, "posts:list", jsonData, 10*time.Minute)

	return postResponses, nil
}

func (ps *PostService) UpdatePost(c *gin.Context, content string, title string, userId uint, postId uint, cfg *config.Config) error {
	post, err := validateUpdate(userId, postId, ps)
	if err != nil {
		return err
	}

	imageURL, err := s3.UploadFile(c, cfg, ps.minio)
	if err != nil {
		return err
	}

	post.Content = content
	post.Title = title
	post.ImageURL = imageURL

	err = ps.postRepo.UpdatePost(post)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("post:%d", postId)
	// Remove the outdated cache
	ps.redis.Del(ctx, cacheKey)
	ps.redis.Publish(ctx, "posts:updates", "update")
	return nil
}

func (ps *PostService) DeletePost(postId uint, userId uint) error {
	post, err := ps.postRepo.GetPostByID(postId)
	if err != nil {
		return errors.New("post not found")
	}
	if post.User.ID != userId {
		return errors.New("user not allowed")
	}

	err = ps.postRepo.DeletePost(postId)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("post:%d", postId)
	log.Println("Delete post:", postId)
	ps.redis.Del(ctx, cacheKey)
	ps.redis.Publish(ctx, "posts:updates", "update")
	return nil
}

func (ps *PostService) GetPostByID(id uint) (*models.Post, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("post:%d", id)

	// Try to get the post from Redis cache
	cachedPost, err := ps.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var post models.Post
		json.Unmarshal([]byte(cachedPost), &post)
		return &post, nil
	}

	// Fetch from DB if not cached
	post, err := ps.postRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// Store in Redis with a TTL of 10 minutes
	jsonData, _ := json.Marshal(post)
	ps.redis.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	return post, nil
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
