package model

import "time"

type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	Count     int
	UserID    uint
	User      User      `gorm:"constraint:OnDelete:CASCADE;"`
	Comments  []Comment `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
