package routes

import (
	"github.com/gin-gonic/gin"
	"go_crud_example/internal/config"
	"go_crud_example/pkg/middleware"
	"gorm.io/gorm"
)

// SetupRoutes подключает все маршруты
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Middleware аутентификации
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Подключаем отдельные роутеры
	SetupAuthRoutes(r, db, cfg)
	SetupPostRoutes(r, db, authMiddleware)
}
