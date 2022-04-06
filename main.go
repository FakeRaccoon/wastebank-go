package main

import (
	"net/http"
	"wastebank/configs"
	"wastebank/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "hello world"}) })
	configs.ConnectDB()

	public := r.Group("/api")

	routes.AuthRoute(public)
	routes.UserRoute(public)
	r.HandleMethodNotAllowed = true
	r.Run()
}
