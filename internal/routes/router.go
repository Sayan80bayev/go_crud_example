package routes

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go_crud_example/internal/config"
	"go_crud_example/pkg/middleware"
	"gorm.io/gorm"
)

// SetupRoutes подключает все маршруты
func SetupRoutes(r *gin.Engine, db *gorm.DB, client *redis.Client, producer *kafka.Producer, cfg *config.Config) {
	// Middleware аутентификации
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Подключаем отдельные роутеры
	SetupAuthRoutes(r, db, cfg)
	SetupPostRoutes(r, db, authMiddleware, cfg)
	SetupCategoryRoutes(r, db, authMiddleware)
	SetupLikeRoutes(r, db, authMiddleware, client, producer)
}
