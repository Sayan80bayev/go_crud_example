package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go_crud_example/internal/repository"
	"log"
)

type CacheListener struct {
	redis    *redis.Client
	postRepo *repository.PostRepository
}

func NewCacheListener(redis *redis.Client, postRepo *repository.PostRepository) *CacheListener {
	return &CacheListener{redis: redis, postRepo: postRepo}
}

func (cl *CacheListener) ListenForPostUpdates() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔥 Паника в ListenForPostUpdates: %v\n", r)
		}
	}()

	ctx := context.Background()
	sub := cl.redis.Subscribe(ctx, "posts:updates")
	ch := sub.Channel()

	log.Println("📡 Подписка на канал posts:updates запущена...")

	for msg := range ch {
		log.Printf("🔄 Получено сообщение: Channel=%s | Payload=%s\n", msg.Channel, msg.Payload)
		cl.redis.Del(ctx, "posts:list")
		UpdateCache(cl.redis, cl.postRepo)
	}
}
