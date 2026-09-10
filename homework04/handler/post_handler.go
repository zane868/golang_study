package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zane868/golang_study/homework04/model"
	"github.com/zane868/golang_study/homework04/service"
	"github.com/zane868/golang_study/homework04/util"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) PublishBlog(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, parseValidationErrors(err))
		return
	}

	post, err := h.postService.CreatePost(req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	util.Success(c, model.PostResponse{
		ID: post.ID,
	})
}
