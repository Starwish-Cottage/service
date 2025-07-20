package admin

import (
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/admin")

	adminGroup.POST("/login", func(c *gin.Context) {
		LoginHandler(c)
	})
}
