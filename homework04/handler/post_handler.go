package handler

import (
	"net/http"
	"strconv"
	"time"

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

func (h *PostHandler) List(c *gin.Context) {
	username := c.GetString("username")

	posts, err := h.postService.List(username)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	response := make([]model.PostResponse, 0, len(posts))
	for _, post := range posts {
		response = append(response, toPostResponse(post))
	}
	util.Success(c, response)
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

	util.Success(c, toPostResponse(post))
}

func (h *PostHandler) Get(c *gin.Context) {
	post, ok := h.getOwnedPost(c)
	if !ok {
		return
	}
	util.Success(c, toPostResponse(post))
}

func (h *PostHandler) Delete(c *gin.Context) {
	post, ok := h.getOwnedPost(c)
	if !ok {
		return
	}
	if err := h.postService.Delete(post.ID); err != nil {
		util.HandleError(c, err)
		return
	}
	util.Success(c, gin.H{"id": post.ID})
}

func (h *PostHandler) Update(c *gin.Context) {
	post, ok := h.getOwnedPost(c)
	if !ok {
		return
	}
	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, parseValidationErrors(err))
		return
	}
	// 以路径中的 ID 为准，防止请求体指定另一篇文章。
	req.PostId = strconv.FormatUint(uint64(post.ID), 10)
	updated, err := h.postService.Update(req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	util.Success(c, toPostResponse(updated))
}

func (h *PostHandler) getOwnedPost(c *gin.Context) (*model.Post, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		util.Error(c, http.StatusBadRequest, "Invalid post ID")
		return nil, false
	}
	userID := c.GetUint("userID")
	if userID == 0 {
		util.Error(c, http.StatusUnauthorized, "Authentication required")
		return nil, false
	}
	post, err := h.postService.Get(uint(id))
	if err != nil {
		util.HandleError(c, err)
		return nil, false
	}
	if post.UserID != userID {
		util.Error(c, http.StatusForbidden, "Access denied")
		return nil, false
	}
	return post, true
}

func toPostResponse(post *model.Post) model.PostResponse {
	const layout = "2006-01-02 15:04:05.000"

	return model.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.In(time.Local).Format(layout),
		UpdatedAt: post.UpdatedAt.In(time.Local).Format(layout),
	}
}
