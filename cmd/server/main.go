package main

import (
	"fmt"
	"go_crud_example/delivery"
	"go_crud_example/internal/config"
	"go_crud_example/internal/models"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db.AutoMigrate(&models.Post{})

	// Инициализация компонентов
	postRepo := repository.NewPostRepository(db)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postHandler := delivery.NewPostHandler(postUsecase)

	// Создаём роутер
	r := gin.Default()
	r.POST("/posts", postHandler.CreatePost)
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPostByID)
	r.PUT("/posts/:id", postHandler.UpdatePost)
	r.DELETE("/posts/:id", postHandler.DeletePost)

	// Запуск сервера
	port := cfg.Port
	fmt.Println("Сервер запущен на порту:", port)
	r.Run(":" + port)
}
