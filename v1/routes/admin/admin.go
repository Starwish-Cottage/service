package admin

import (
	"github.com/Starwish-Cottage/service/v1/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")

	adminGroup.POST("/login", func(c *gin.Context) {
		LoginHandler(c)
	})

	// Serve static files
	adminGroup.Static("/images", "./scripts/src_imgs")

	protected := adminGroup.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/verify-session", func(c *gin.Context) {
			VerifySessionHandler(c)
		})

		protected.POST("/upload", func(c *gin.Context) {
			UploadImageHandler(c)
		})
	}
}
