package controllers

import (
	"context"
	"net/http"
	"time"
	"wastebank/models"
	"wastebank/responses"
	"wastebank/services"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input models.LoginInput
		var err error
		defer cancel()

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(500, responses.LoginResponse{Status: 500, Data: err.Error()})
			return
		}

		u := models.UserWithPassword{}
		ui := models.LoginInput{}

		ui.Username = input.Username
		ui.Password = input.Password

		err = UserCollection.FindOne(ctx, bson.M{"username": ui.Username}).Decode(&u)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}

		err = services.VerifyPassword(ui.Password, u.Password)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password is incorrect."})
			return
		}

		jwt := services.Jwt{}

		token, err := jwt.GenerateAllToken(u)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot generate token."})
			return
		}

		c.JSON(http.StatusOK, token)

	}
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := models.RefreshToken{}
		c.ShouldBindBodyWith(&token, binding.JSON)

		jwt := services.Jwt{}
		user, errRefresh := jwt.ValidateRefreshToken(token)
		if errRefresh != nil {
			c.JSON(401, gin.H{"message": "invalid refresh token"})
			return
		}

		accessToken, err := jwt.GenerateToken(user)

		if err != nil {
			c.JSON(401, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, accessToken)

	}
}
