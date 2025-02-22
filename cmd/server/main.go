package main

import (
	"fmt"
	"go_crud_example/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	"go_crud_example/internal/config"
	"go_crud_example/internal/routes"
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
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Автоматическая миграция
	db.AutoMigrate(&models.User{}, &models.Post{})

	// Создаём роутер и передаём зависимости
	r := gin.Default()
	routes.SetupRoutes(r, db, cfg)

	// Запуск сервера
	fmt.Println("Сервер запущен на порту:", cfg.Port)
	r.Run(":" + cfg.Port)
}
