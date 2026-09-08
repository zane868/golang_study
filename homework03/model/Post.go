package model

import (
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

type Post struct {
	ID            uint `gorm:"primaryKey"`
	Title         string
	Content       string
	Count         int
	UserID        uint
	User          User      `gorm:"constraint:OnDelete:CASCADE;"`
	Comments      []Comment `gorm:"constraint:OnDelete:CASCADE;"`
	CommentStatus string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (p *Post) BeforeCreate(db *gorm.DB) error {
	p.Count = utf8.RuneCountInString(p.Content)
	return nil
}
