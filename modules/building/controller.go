package buildingM

import (
	"github.com/gin-gonic/gin"
	modal "github.com/windyzoe/study-house/modules/building/model"
	"github.com/windyzoe/study-house/util"
)

func SetupRouter(router *gin.Engine) {
	group := router.Group("/building")
	{
		group.GET("/list", list)
	}
}

func list(c *gin.Context) {
	var req modal.BuildingListReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorRes(c, 10015, err)
		return
	}
	pageData := GetBuildingList(req.PageNumber, req.PageSize)
	util.SuccessRes(c, pageData)
}
