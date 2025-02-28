package routes

import (
	"github.com/gin-gonic/gin"
	"go_crud_example/internal/delivery"
	"go_crud_example/internal/repository"
	"go_crud_example/internal/service"
	"gorm.io/gorm"
)

func SetupCategoryRoutes(router *gin.Engine, db *gorm.DB, authMiddleware gin.HandlerFunc) {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := delivery.NewCategoryHandler(categoryService)

	router.GET("/category", categoryHandler.ListCategory)

	categoryGroup := router.Group("/category", authMiddleware)
	{
		categoryGroup.POST("/", categoryHandler.CreateCategory)
		categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)
	}

}
