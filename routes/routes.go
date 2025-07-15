package routes

import (
	"github.com/Starwish-Cottage/service/routes/admin"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	admin.SetupAdminRoutes(router)
}
