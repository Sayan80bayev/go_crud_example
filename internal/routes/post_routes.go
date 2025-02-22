package routes

import (
	"github.com/gin-gonic/gin"
	"go_crud_example/internal/delivery"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/service"
	"gorm.io/gorm"
)

// SetupPostRoutes настраивает маршруты для работы с постами
func SetupPostRoutes(r *gin.Engine, db *gorm.DB, authMiddleware gin.HandlerFunc) {
	postRepo := repository.NewPostRepository(db)
	postUsecase := service.NewPostService(postRepo)
	postHandler := delivery.NewPostHandler(postUsecase)

	// Открытые роуты
	r.GET("/posts", postHandler.GetPosts)
	r.GET("/posts/:id", postHandler.GetPostByID)

	// Защищённые роуты (требуется авторизация)
	postRoutes := r.Group("/posts", authMiddleware)
	{
		postRoutes.POST("/", postHandler.CreatePost)
	}
}
