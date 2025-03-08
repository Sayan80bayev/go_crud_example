package mappers

import (
	"go_crud_example/internal/models"
	"go_crud_example/internal/response"
)

func MapPostsToResponse(posts []models.Post) []response.PostResponse {
	var postResponses []response.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, response.PostResponse{
			ID:    post.ID,
			Title: post.Title,
			Author: response.UserResponse{
				ID:       post.User.ID,
				Username: post.User.Username,
			},
			Category: response.CategoryResponse{
				Id:   post.Category.ID,
				Name: post.Category.Name,
			},
			ImageURL:  &post.ImageURL,
			LikeCount: post.LikeCount,
		})
	}
	return postResponses
}
