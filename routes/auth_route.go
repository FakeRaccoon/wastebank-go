package routes

import (
	"wastebank/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.RouterGroup) {
	r.POST("login", controllers.Login())
	r.POST("refresh", controllers.RefreshToken())
}
