package cache

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"go_crud_example/internal/mappers"
	"go_crud_example/internal/repository"
)

func UpdateCache(redis *redis.Client, postRepo *repository.PostRepository) {
	ctx := context.Background()

	posts, err := postRepo.GetPosts()
	if err != nil {
		log.Println("Ошибка загрузки постов:", err)
		return
	}

	// Конвертируем в JSON
	postResponses := mappers.MapPostsToResponse(posts)
	jsonData, _ := json.Marshal(postResponses)

	// Обновляем кэш
	redis.Set(ctx, "posts:list", jsonData, 0) // 0 – без истечения времени
	log.Println("✅ Кэш обновлён!")
}
