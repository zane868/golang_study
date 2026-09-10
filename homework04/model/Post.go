package model

import (
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Count     int
	UserID    uint
	User      User `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePostRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

type PostResponse struct {
	ID uint `json:"id"`
}

func (p *Post) BeforeCreate(db *gorm.DB) error {
	p.Count = utf8.RuneCountInString(p.Content)
	return nil
}
