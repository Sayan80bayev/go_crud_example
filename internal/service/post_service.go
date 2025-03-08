package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"

	"go_crud_example/internal/config"
	"go_crud_example/internal/mappers"
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
	mapper   *mappers.PostMapper
}

// Конструктор сервиса PostService
func NewPostService(postRepo *repository.PostRepository, minioClient *minio.Client, redis *redis.Client, producer *kafka.Producer) *PostService {
	return &PostService{
		postRepo: postRepo,
		minio:    minioClient,
		redis:    redis,
		producer: producer,
		mapper:   mappers.NewPostMapper(),
	}
}

// Создание поста
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

// Получение списка постов
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

	postResponses := ps.mapper.MapEach(posts)
	jsonData, _ := json.Marshal(postResponses)
	ps.redis.Set(ctx, "posts:list", jsonData, 10*time.Minute)

	return postResponses, nil
}

// Получение поста по ID
func (ps *PostService) GetPostByID(id uint) (*response.PostResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("post:%d", id)

	// Попытка получить пост из кэша
	cachedPost, err := ps.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var post response.PostResponse
		json.Unmarshal([]byte(cachedPost), &post)
		return &post, nil
	}

	// Если в кэше нет, берем из БД
	post, err := ps.postRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	// Маппинг в DTO
	resp := ps.mapper.Map(*post)
	jsonData, _ := json.Marshal(resp)
	ps.redis.Set(ctx, cacheKey, jsonData, 10*time.Minute)

	return &resp, nil
}

// Обновление поста
func (ps *PostService) UpdatePost(c *gin.Context, content string, title string, userId uint, postId uint, cfg *config.Config) error {
	post, err := validatePermission(userId, postId, ps)
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
	ps.redis.Del(ctx, cacheKey)
	ps.redis.Publish(ctx, "posts:updates", "update")

	return nil
}

// Удаление поста
func (ps *PostService) DeletePost(postId uint, userId uint) error {
	_, err := validatePermission(userId, postId, ps)

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

// Валидация перед обновлением
func validatePermission(userId uint, postId uint, ps *PostService) (*models.Post, error) {
	post, err := ps.postRepo.GetPostByID(postId)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.User.ID != userId {
		return nil, errors.New("user not allowed")
	}

	return post, nil
}
