package models

type Like struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id" gorm:"not null;index"`
	PostID uint `json:"post_id" gorm:"not null;index"`
}
