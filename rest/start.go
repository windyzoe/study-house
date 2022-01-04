package rest

import (
	"github.com/gin-gonic/gin"
	buildingM "github.com/windyzoe/study-house/modules/building"
	houseM "github.com/windyzoe/study-house/modules/house"
	schoolM "github.com/windyzoe/study-house/modules/school"
	userM "github.com/windyzoe/study-house/modules/user"
	"github.com/windyzoe/study-house/util"
)

var Router *gin.Engine

func Start() {
	Router = gin.Default()
	Router.Use(userM.Auth())
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	setupRouter()
	Router.Run(":" + util.Configs.Server.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupRouter() {
	userM.SetupRouter(Router)
	buildingM.SetupRouter(Router)
	schoolM.SetupRouter(Router)
	houseM.SetupRouter(Router)
}
