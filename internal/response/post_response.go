package response

type PostResponse struct {
	ID       uint             `json:"id"`
	Title    string           `json:"title"`
	Author   UserResponse     `json:"author"`
	Category CategoryResponse `json:"category"`
}
