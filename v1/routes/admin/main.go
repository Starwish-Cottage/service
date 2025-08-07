package admin

import (
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")

	adminGroup.POST("/login", func(c *gin.Context) {
		LoginHandler(c)
	})

	adminGroup.POST("/verify-session", func(c *gin.Context) {
		VerifySessionHandler(c)
	})

	// Serve static files
	adminGroup.Static("/images", "./scripts/src_imgs")
	adminGroup.POST("/upload", func(c *gin.Context) {
		UploadImageHandler(c)
	})
}
