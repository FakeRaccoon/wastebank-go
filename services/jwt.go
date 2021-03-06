package services

import (
	"errors"
	"fmt"
	"os"
	"time"
	"wastebank/models"

	"github.com/golang-jwt/jwt"
)

type Jwt struct{}

func (j Jwt) GenerateAllToken(user models.UserWithPassword) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["user_id"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt := models.Token{}

	jwt.AccessToken, err = token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return jwt, err
	}

	return j.createRefreshToken(jwt)
}

func (Jwt) GenerateToken(user models.UserWithPassword) (models.AccessToken, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["user_id"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt := models.AccessToken{}

	jwt.AccessToken, err = token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return jwt, err
	}

	return jwt, nil
}

func (Jwt) ValidateToken(accessToken string) (models.UserWithPassword, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})

	user := models.UserWithPassword{}
	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user.Username = payload["user_id"].(string)

		return user, nil
	}

	return user, errors.New("invalid token")
}

func (j Jwt) ValidateRefreshToken(model models.RefreshToken) (models.UserWithPassword, error) {
	token, err := jwt.Parse(model.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})

	user := models.UserWithPassword{}
	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return user, errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err = parser.ParseUnverified(payload["token"].(string), claims)
	if err != nil {
		return user, err
	}

	payload, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid token")
	}

	user.Username = payload["user_id"].(string)

	return user, nil
}

func (Jwt) createRefreshToken(token models.Token) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["token"] = token.AccessToken
	claims["exp"] = time.Now().Add(time.Hour * 24 * 90).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return token, err
	}

	return token, nil
}
