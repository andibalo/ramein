package util

import (
	"github.com/andibalo/ramein/core/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"log"
)

func GenerateToken(user *model.User) (tokenString string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.FirstName + " " + user.LastName,
		"email": user.Email,
		"phone": user.Phone,
		"role":  user.Role,
	})

	tokenString, err = token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Println(err)

		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to sign JWT token")
	}

	return tokenString, nil
}
