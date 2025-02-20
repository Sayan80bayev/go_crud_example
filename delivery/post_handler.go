package delivery

import (
	"github.com/gin-gonic/gin"
	"go_crud_example/internal/models"
	"go_crud_example/internal/usecase"
	"net/http"
	"strconv"
)

type PostHandler struct {
	usecase *usecase.PostUsecase
}

func NewPostHandler(uc *usecase.PostUsecase) *PostHandler {
	return &PostHandler{usecase: uc}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.usecase.CreatePost(&post)
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.usecase.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := h.usecase.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.ID = uint(id)
	h.usecase.UpdatePost(&post)
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.usecase.DeletePost(uint(id))
	c.JSON(http.StatusOK, gin.H{"message": "Пост удалён"})
}
