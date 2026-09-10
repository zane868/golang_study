package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/service"
	"github.com/zane868/golang_study/homework04/util"
)

type UserHandler struct {
	userService *service.UserService
	jwtSecret   []byte
}

func NewUserHandler(userService *service.UserService, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, parseValidationErrors(err))
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, parseValidationErrors(err))
		return
	}

	user, err := h.userService.Authenticate(req.Username, req.Password)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	token, err := util.GenerateToken(h.jwtSecret, user.ID, user.Username)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.Header("Authorization", token)
	c.SetCookie("Authorization", token, 3600, "/", "", false, true)
	util.Success(c, gin.H{
		"authorization": token,
		"user": model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

func parseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	// 简化处理，实际应该解析 binding 错误
	errors["general"] = err.Error()
	return errors
}
