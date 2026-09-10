package service

import (
	"errors"

	"github.com/zane868/golang_study/homework04/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/zane868/golang_study/homework04/util"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(req model.CreateUserRequest) (*model.User, error) {

	// 检查用户名是否已存在
	if err := s.CheckNameExists(req.Username); err != nil {
		return nil, err
	}

	// 检查邮箱是否已存在
	if err := s.CheckEmailExists(req.Email); err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CheckNameExists(name string) error {
	return s.CheckExists("username = ?", name)
}

func (s *UserService) CheckEmailExists(email string) error {
	return s.CheckExists("email = ?", email)
}

func (s *UserService) CheckExists(query interface{}, args ...interface{}) error {
	var count int64
	err := s.db.Model(&model.User{}).Where(query, args).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return util.NewAppError(409, "Username already exists")

	}
	return nil
}

func (s *UserService) Authenticate(username, password string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.NewAppError(401, "Invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, util.NewAppError(401, "Invalid credentials")
	}

	return &user, nil
}
