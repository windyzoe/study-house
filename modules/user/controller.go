package userM

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/windyzoe/study-house/util"
)

type LoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func SetupRouter(router *gin.Engine) {
	v1 := router.Group("/login")
	{
		v1.POST("", Login)
	}
	router.POST("/validateToken", ValidateToken)
}

func Login(c *gin.Context) {
	var current LoginReq
	if err := c.ShouldBind(&current); err != nil {
		util.ErrorRes(c, 10015, err)
		return
	}
	user, err := GetUser(current.Name)
	if err != nil {
		log.Println(err)
		util.ErrorRes(c, 10016, err)
		return
	}
	if user["Password"] != current.Password {
		util.ErrorRes(c, 10017, err)
		return
	}
	item := GenerateToken()
	util.SuccessRes(c, item.Token)
}

func ValidateToken(c *gin.Context) {
	util.SuccessRes(c, nil)
}
