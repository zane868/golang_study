package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code:  code,
		Msg:   msg,
		Error: msg,
	})
}

func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:  422,
		Msg:   "validation failed",
		Error: errors,
	})
}
