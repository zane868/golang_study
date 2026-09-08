package model

import (
	"time"

	"gorm.io/gorm"
)

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

func (c *Comment) AfterDelete(db *gorm.DB) error {
	if c.PostID <= 0 {
		return nil
	}
	var count int64
	db.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)
	if count <= 0 {
		db.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_status", "无评论")
	}
	return nil
}
