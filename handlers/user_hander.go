package handlers

import (
	"context"
	"user_login/models"
	"user_login/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

func GetUser(c *fiber.Ctx) error {
	accessToken := c.Get("accessToken")
	tokenAcc, err := utils.ParseToken(accessToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	if tokenAcc.Valid {
		claims := tokenAcc.Claims.(jwt.MapClaims)
		username := claims["username"].(string)
		email := claims["email"].(string)

		user := models.User{}
		filter := bson.M{"email": email, "username": username}
		err := userCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		response := map[string]string{
			"username": user.Username,
			"email":    user.Email,
			"created":  user.CreatedAt.String(),
		}

		return c.Status(fiber.StatusOK).JSON(response)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token is invalid"})
	}

}
