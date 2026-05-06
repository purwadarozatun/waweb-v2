package middleware

import (
	"strings"
	"waweb-v2/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookie
		tokenString := c.Cookies("token")

		if tokenString == "" {
			// Check Authorization header
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return c.Redirect("/login")
		}

		// Parse token
		cfg := config.LoadConfig()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.ClearCookie("token")
			return c.Redirect("/login")
		}

		// Get claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("userID", claims["user_id"])
			c.Locals("email", claims["email"])
		}

		return c.Next()
	}
}
