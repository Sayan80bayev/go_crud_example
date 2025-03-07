package main

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go_crud_example/internal/config"
	"go_crud_example/internal/models"
	"go_crud_example/internal/routes"
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

	db.AutoMigrate(&models.User{}, &models.Post{})

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass, // Оставьте пустым, если пароль не нужен
		DB:       0,             // Используем стандартную базу
	})

	// Проверка соединения с Redis
	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Ошибка подключения к Redis:", err)
	}

	// Подключение к Kafka
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.KafkaBrokers[0]})
	if err != nil {
		log.Fatal("Ошибка подключения к Kafka:", err)
	}
	defer producer.Close()

	log.Println("✅ Подключение к базе данных, Redis и Kafka успешно установлено")

	// Создаём роутер и передаём зависимости
	r := gin.Default()
	routes.SetupRoutes(r, db, redisClient, producer, cfg)

	// Запуск сервера
	r.Run(":" + cfg.Port)
}
