package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"
	"wastebank/configs"
	"wastebank/models"
	"wastebank/responses"
	"wastebank/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

var validate = validator.New()

func GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		bearerToken := c.Request.Header["Authorization"]
		token := strings.Split(bearerToken[0], " ")
		defer cancel()
		var u models.User

		jwt := services.Jwt{}
		user, errToken := jwt.ValidateToken(token[1])

		if errToken != nil {
			c.JSON(401, responses.UserResponse{Status: 401, Data: errToken.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&u)

		if err != nil {
			c.JSON(401, responses.UserResponse{Status: 401, Data: errToken.Error()})
			return
		}

		c.JSON(200, responses.UserResponse{Status: 200, Data: u})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("id")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Data: err.Error()})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Data: user})
	}
}

func GetAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err := UserCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(500, responses.UserResponse{Status: 500, Data: err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(500, responses.UserResponse{Status: 500, Data: err.Error()})
			}

			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: 200, Data: users})
	}
}
