package models

type Post struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`           // ID автора поста
	User    User   `gorm:"foreignKey:UserID"` // Связь с автором
}
