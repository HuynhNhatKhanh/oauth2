// handlers/user_handler.go
package handlers

import (
	"user_login/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	// Lấy thông tin user từ token (accessToken)
	accessToken := c.Get("Authorization")

	tokenAcc, err := utils.ParseToken(accessToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	if tokenAcc.Valid {
		claims := tokenAcc.Claims.(jwt.MapClaims)
		username := claims["username"].(string)
		email := claims["email"].(string)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"username": username, "email": email, "message": "User info retrieved successfully"})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token is invalid"})
	}

}
