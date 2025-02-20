package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_crud_example/delivery"
	"go_crud_example/internal/config"
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/usecase"
	"go_crud_example/pkg/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}

	// Миграция
	db.AutoMigrate(&models.Post{}, &models.User{})

	// Инициализация компонентов
	postRepo := repository.NewPostRepository(db)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postHandler := delivery.NewPostHandler(postUsecase)

	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	authHandler := delivery.NewAuthHandler(authUsecase)

	// Используем единый ключ для всех компонентов
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Создаём роутер
	r := gin.Default()
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPostByID)

	postRoutes := r.Group("/posts", authMiddleware)
	{
		postRoutes.POST("/", postHandler.CreatePost)
	}

	r.POST("/register", authHandler.Register)
	r.POST("/auth", authHandler.Login)

	// Запуск сервера
	fmt.Println("Сервер запущен на порту:", cfg.Port)
	r.Run(":" + cfg.Port)
}
