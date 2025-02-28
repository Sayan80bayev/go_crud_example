package models

type Post struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	CategoryID uint     `json:"category_id gorm: "not null default -1"`
	Category   Category `gorm:"foreignKey:CategoryID"` // Связь с категорией
	UserID     uint     `json:"user_id"`               // ID автора поста
	User       User     `gorm:"foreignKey:UserID"`     // Связь с автором
}
