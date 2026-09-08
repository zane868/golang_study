package model

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	Post      Post `gorm:"constraint:OnDelete:CASCADE;"`
	UserId    uint
	User      User `gorm:"constraint:OnDelete:CASCADE;"`
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
