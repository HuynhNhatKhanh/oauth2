// handlers/user_handler.go
package handlers

import (
	"user_login/utils"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	// Lấy thông tin user từ token (accessToken)
	accessToken := c.Get("Authorization")
	userID, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Tìm kiếm user trong database và trả về thông tin
	// ...

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"userID": userID, "message": "User info retrieved successfully"})
}

// Các hàm khác để thực hiện các thao tác khác với người dùng như cập nhật thông tin, đổi mật khẩu, v.v.
