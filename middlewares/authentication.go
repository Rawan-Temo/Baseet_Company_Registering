package middlewares

import (
	"errors"
	"slices"
	"strings"

	"github.com/Rawan-Temo/Baseet_Company_Registering.git/dtos"
	"github.com/Rawan-Temo/Baseet_Company_Registering.git/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuthenticated(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		// No Bearer prefix
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return utils.JwtSecret, nil
	})

	if err != nil {
		// Check if the error is token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Token expired",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	currentUser := dtos.UserTokenClaim{
		UserID: uint(claims["user_id"].(float64)),
		Role:   claims["role"].(string),
		Exp:    int64(claims["exp"].(float64)),
		Iat:    int64(claims["iat"].(float64)),
	}
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Unauthorized",
		})
	}

	c.Locals("currentUser", currentUser)
	return c.Next()
}
func AllowedTo(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser")
		if currentUser == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "fail",
				"message": "Unauthorized",
			})
		}

		claims := currentUser.(jwt.MapClaims)
		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Forbidden",
			})
		}

		// Check if role is allowed
		authorized := slices.Contains(roles, role)

		if !authorized {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Forbidden",
			})
		}

		return c.Next()
	}
}
