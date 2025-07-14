package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/admin")

	adminGroup.GET("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"login": "admin",
		})
	})
}
