package delivery

import (
	"github.com/gin-gonic/gin"
	"go_crud_example/internal/service"
	"net/http"
	"strconv"
)

type LikeHandler struct {
	LikeService service.LikeService
}

func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{likeService}
}

func (h *LikeHandler) CreateLike(c *gin.Context) {
	// Extract user ID and safely convert it to uint
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint, ok := userID.(uint) // Often Gin stores JSON numbers as float64
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Extract post ID from query and convert to uint
	postIDStr := c.Query("post")
	postID, err := strconv.ParseUint(postIDStr, 10, 64) // Convert to uint64 first
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	postIDUint := uint(postID) // Convert to uint

	// Call the service method
	if err := h.LikeService.LikePost(userIDUint, postIDUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"": "Like created successfully"})
}

func (h *LikeHandler) DeleteLike(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	likeIDstr := c.Query("like")
	likeID, err := strconv.ParseUint(likeIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid like ID"})
		return
	}
	likeIDUint := uint(likeID)

	if err := h.LikeService.UnlikePost(likeIDUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"": "Like deleted successfully"})
}
