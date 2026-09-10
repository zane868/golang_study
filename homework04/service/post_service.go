package service

import (
	"github.com/zane868/golang_study/homework04/model"

	"gorm.io/gorm"
)

type PostService struct {
	db          *gorm.DB
	userService *UserService
}

func NewPostService(u *UserService, db *gorm.DB) *PostService {
	return &PostService{db: db, userService: u}
}

func (p *PostService) CreatePost(req model.CreatePostRequest) (*model.Post, error) {

	user, err := p.userService.GetUser(req.Username)
	if err != nil {
		return nil, err
	}

	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  user.ID,
	}

	err = p.db.Create(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
