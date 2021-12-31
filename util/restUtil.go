package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessRes(c *gin.Context, data interface{}) {
	var success ResponseData
	success.Code = 0
	success.Msg = "success"
	success.Data = data
	c.JSON(http.StatusOK, success)
	c.Abort()
}

func ErrorRes(c *gin.Context, status int, err error) {
	var errorData ResponseData
	errorData.Code = status
	errorData.Msg = err.Error()
	errorData.Data = nil
	c.JSON(http.StatusOK, errorData)
	c.Abort()
}

func ErrorWithDataRes(c *gin.Context, status int, err error, data interface{}) {
	var errorData ResponseData
	errorData.Code = status
	errorData.Msg = err.Error()
	errorData.Data = data
	c.JSON(http.StatusOK, errorData)
}
