package service

import (
	"errors"

	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/util"
	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (c *CommentService) List(postID uint) ([]model.Comment, error) {
	var post model.Post
	if err := c.db.First(&post, postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewAppError(404, "Post not found")
		}
		return nil, err
	}
	comments := make([]model.Comment, 0)
	err := c.db.Preload("User").Where("post_id = ?", postID).Order("id ASC").Find(&comments).Error
	return comments, err
}

func (c *CommentService) Create(req model.CreateCommentRequest) (*model.Comment, error) {
	var post model.Post
	if err := c.db.First(&post, req.PostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewAppError(404, "Post not found")
		}
		return nil, err
	}
	comment := model.Comment{
		PostID:  req.PostID,
		UserID:  req.UserId,
		Content: req.Content,
	}
	err := c.db.Create(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *CommentService) Delete(req model.DeleteCommenctRequest) error {
	err := c.db.Delete(&model.Comment{}, req.CommentID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *CommentService) Get(commentId uint) (*model.Comment, error) {
	var comment model.Comment
	if err := p.db.First(&comment, commentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewAppError(404, "Comment not found")
		}
		return nil, err
	}
	return &comment, nil
}
