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

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) List(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Query("postId"), 10, strconv.IntSize)
	if err != nil || postID == 0 {
		util.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}
	comments, err := h.commentService.List(uint(postID))
	if err != nil {
		util.HandleError(c, err)
		return
	}
	response := make([]model.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		response = append(response, toCommentResponse(comment))
	}
	util.Success(c, response)
}

func toCommentResponse(comment model.Comment) model.CommentResponse {
	return model.CommentResponse{
		ID: comment.ID, PostID: comment.PostID, UserID: comment.UserID,
		Username: comment.User.Username, Content: comment.Content,
		CreatedAt: comment.CreatedAt.In(time.Local).Format("2006-01-02 15:04:05.000"),
	}
}

func (h *CommentHandler) Comment(c *gin.Context) {

	userId := c.GetUint("userID")

	req := model.CreateCommentRequest{
		UserId: userId,
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ValidationError(c, parseValidationErrors(err))
		return
	}

	req.UserId = userId
	comment, err := h.commentService.Create(req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	util.Success(c, gin.H{"id": comment.ID, "post_id": comment.PostID, "user_id": comment.UserID, "content": comment.Content})

}

func (h *CommentHandler) Delete(c *gin.Context) {
	comment, ok := h.getOwnedComment(c)
	if !ok {
		return
	}
	err := h.commentService.Delete(model.DeleteCommenctRequest{CommentID: comment.ID})
	if err != nil {
		util.HandleError(c, err)
		return
	}
	util.Success(c, nil)

}

func (h *CommentHandler) getOwnedComment(c *gin.Context) (*model.Comment, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		util.Error(c, http.StatusBadRequest, "Invalid comment ID")
		return nil, false
	}
	userID := c.GetUint("userID")
	if userID == 0 {
		util.Error(c, http.StatusUnauthorized, "Authentication required")
		return nil, false
	}
	comment, err := h.commentService.Get(uint(id))
	if err != nil {
		util.HandleError(c, err)
		return nil, false
	}
	if comment.UserID != userID {
		util.Error(c, http.StatusForbidden, "Access denied")
		return nil, false
	}
	return comment, true
}
