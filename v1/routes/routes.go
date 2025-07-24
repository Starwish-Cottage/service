package routes

import (
	"github.com/Starwish-Cottage/service/v1/routes/admin"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	admin.SetupAdminRoutes(v1)
}
