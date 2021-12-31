package houseM

import (
	"github.com/gin-gonic/gin"
	"github.com/windyzoe/study-house/modules/house/modal"
	"github.com/windyzoe/study-house/util"
)

func SetupRouter(router *gin.Engine) {
	group := router.Group("/house")
	{
		group.GET("/list", list)
	}
}

func list(c *gin.Context) {
	var req modal.HouseListReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorRes(c, 10015, err)
		return
	}
	pageData := GetHouseList(req.PageNumber, req.PageSize)
	util.SuccessRes(c, pageData)
}
