package models

type Like struct {
	ID     uint `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uint `json:"user_id" gorm:"not null;index"`
	PostID uint `json:"post_id" gorm:"not null;index"`
}
