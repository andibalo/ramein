package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"strings"
)

// TokenClaims : struct for validate token claims
type TokenClaims struct {
	ID    interface{} `json:"id"`
	Email string      `json:"email"`
	Name  string      `json:"name"`
	Phone string      `json:"phone"`
	Role  string      `json:"role"`
	jwt.RegisteredClaims
}

// contextClaimKey key value store/get token on context
const contextClaimKey = "ctx.mw.auth.claim"

func GetToken(c *fiber.Ctx) string {
	authorization := c.Get("Authorization")
	tokens := strings.Split(authorization, "Bearer ")

	return tokens[1]
}

func CheckIsAuthenticated(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid token"})
	}

	if tokenString == viper.Get("STATIC_JWT") {
		c.Set("x-user-email", "admin@ramein.com")
		c.Locals(contextClaimKey, &TokenClaims{
			Email: "admin@ramein.com",
			Name:  "superadmin",
		})

		return c.Next()
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(viper.GetString("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(TokenClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})
	}

	c.Set("x-user-email", claims.Email)
	c.Locals(contextClaimKey, claims)

	return c.Next()
}
