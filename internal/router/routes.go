package router

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
	SetupManualCutterRoutes(r)
}
