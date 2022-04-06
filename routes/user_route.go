package routes

import (
	"wastebank/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.RouterGroup) {
	r.GET("/users", controllers.GetAllUser())
	r.GET("/users/current", controllers.GetCurrentUser())
	r.GET("/users/:id", controllers.GetUser())
}
