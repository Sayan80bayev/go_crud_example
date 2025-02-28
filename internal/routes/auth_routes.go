package routes

import (
	"github.com/gin-gonic/gin"
	"go_crud_example/internal/config"
	"go_crud_example/internal/delivery"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/service"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := delivery.NewAuthHandler(authService)

	// Роуты для аутентификации
	r.POST("/register", authHandler.Register)
	r.POST("/auth", authHandler.Login)
}
