package service

import (
	"errors"
	"strconv"
	"unicode/utf8"

	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/util"

	"gorm.io/gorm"
)

type PostService struct {
	db          *gorm.DB
	userService *UserService
}

func NewPostService(u *UserService, db *gorm.DB) *PostService {
	return &PostService{db: db, userService: u}
}

func (p *PostService) List(username string) ([]*model.Post, error) {
	user, err := p.userService.GetUser(username)
	if err != nil {
		return nil, err
	}
	posts := make([]*model.Post, 0)
	if err := p.db.
		Where("user_id = ?", user.ID).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostService) Get(postId uint) (*model.Post, error) {
	var post model.Post
	if err := p.db.First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewAppError(404, "Post not found")
		}
		return nil, err
	}

	return &post, nil
}

func (p *PostService) Delete(postId uint) error {
	post, err := p.Get(postId)
	if err != nil {
		return err
	}

	return p.db.Delete(post).Error
}

func (p *PostService) Update(req model.UpdatePostRequest) (*model.Post, error) {
	postID, err := strconv.ParseUint(req.PostId, 10, 64)
	if err != nil || postID == 0 {
		return nil, util.NewAppError(400, "Invalid post ID")
	}

	post, err := p.Get(uint(postID))
	if err != nil {
		return nil, err
	}

	post.Title = req.Title
	post.Content = req.Content
	post.Count = utf8.RuneCountInString(req.Content)

	if err := p.db.Save(post).Error; err != nil {
		return nil, err
	}

	return post, nil
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
