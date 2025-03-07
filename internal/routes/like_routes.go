package routes

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go_crud_example/internal/delivery"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/service"
	"gorm.io/gorm"
)

func SetupLikeRoutes(r *gin.Engine, db *gorm.DB, authMiddleware gin.HandlerFunc, redisClient *redis.Client, kafkaProd *kafka.Producer) {

	likeRepository := repository.NewLikeRepository(db)
	likeService := service.NewLikeService(likeRepository)
	likeHandler := delivery.NewLikeHandler(likeService)

	likeRoutes := r.Group("/like", authMiddleware)
	{
		likeRoutes.POST("/", likeHandler.CreateLike)
		likeRoutes.DELETE("/", likeHandler.DeleteLike)
	}
}
