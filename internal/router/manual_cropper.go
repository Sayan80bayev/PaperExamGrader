package router

import (
	"github.com/gin-gonic/gin"
)

func SetupManualCutterRoutes(r *gin.Engine) {

	r.POST("/api/v1/crop/manual", func(c *gin.Context) {

	})

}
