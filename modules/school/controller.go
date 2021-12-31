package schoolM

import (
	"github.com/gin-gonic/gin"
	"github.com/windyzoe/study-house/modules/school/modal"
	"github.com/windyzoe/study-house/util"
)

func SetupRouter(router *gin.Engine) {
	group := router.Group("/school")
	{
		group.GET("/list", list)
	}
}

func list(c *gin.Context) {
	var req modal.SchoolListReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorRes(c, 10015, err)
		return
	}
	pageData := GetSchoolList(req.PageNumber, req.PageSize)
	util.SuccessRes(c, pageData)
}
