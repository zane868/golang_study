package model

import (
	"time"
)

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	Post      Post `gorm:"constraint:OnDelete:CASCADE;"`
	UserID    uint
	User      User `gorm:"constraint:OnDelete:CASCADE;"`
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateCommentRequest struct {
	PostID  uint   `json:"postId" binding:"required"`
	UserId  uint   `json:"-"`
	Content string `json:"content" binding:"required"`
}

type DeleteCommenctRequest struct {
	CommentID uint `json:"commentId"`
}

type CommentResponse struct {
	ID        uint   `json:"id"`
	PostID    uint   `json:"post_id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
