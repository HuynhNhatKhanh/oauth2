// handlers/auth_handler.go
package handlers

import (
	"context"
	"fmt"
	"time"
	"user_login/config"
	"user_login/models"
	"user_login/utils"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func Register(c *fiber.Ctx) error {

	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if err := validate.Struct(newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Check if user already exists
	if CheckUserExist(c, newUser.Email) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User already exists"})
	}
	otp := utils.GenerateOTP()
	newUser.OTP = otp
	newUser.Password = utils.HashPassword(newUser.Password)
	newUser.IsVerified = false
	newUser.CreatedAt = time.Now()

	_, err := userCollection.InsertOne(context.Background(), newUser)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving user"})
	}

	//Send email verification
	errMail := utils.SendOTP(newUser.Email, otp)
	if errMail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP"})
	}
	return c.Redirect("/verify", fiber.StatusSeeOther)
	// return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Send otp successfully"})
}

func VerifyEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	code := c.Query("otp")

	// Retrieve user from the database
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return err
	}

	// Check if the user's email is already verified
	if user.IsVerified {
		return c.JSON(fiber.Map{"message": "Email already verified"})
	}

	if code != user.OTP {
		return c.JSON(fiber.Map{"message": "Invalid verification code"})
	}

	// Update the user's verification status in the database
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"is_verified": true}}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Email verified successfully"})
}

func Login(c *fiber.Ctx) error {

	// Get user from request body
	var existingUser models.User
	if err := c.BodyParser(&existingUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Check and find email user in database
	filter := bson.M{"email": existingUser.Email}
	result := userCollection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		fmt.Println(result.Err())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Decode user from database
	var userFromDB models.User
	if err := result.Decode(&userFromDB); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	// Check password
	if !utils.ComparePasswordHash(userFromDB.Password, existingUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Incorrect password"})
	}

	// Check if user is verified
	if !userFromDB.IsVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email not verified"})
	}

	// Check if user is verified login
	if !userFromDB.IsVerifiedLogin {
		otp := utils.GenerateOTP()
		errMail := utils.SendOTP(userFromDB.Email, otp)
		if errMail != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP"})
		}
		userCollection.UpdateOne(context.Background(), bson.M{"email": userFromDB.Email}, bson.M{"$set": bson.M{"otp_login": otp}})

		return c.Redirect("/profile", fiber.StatusSeeOther)
	}

	// Tạo accessToken và refreshToken
	// accessToken, refreshToken, err := utils.CreateTokens(existingUser.ID)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating token"})
	// }

	// // Lưu refreshToken vào database
	// existingUser.RefreshToken = refreshToken
	// err = existingUser.UpdateRefreshToken(userCollection)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving refresh token"})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successfully"})
}

func VerifyLogin(c *fiber.Ctx) error {
	email := c.Query("email")
	code := c.Query("otp_login")

	// Retrieve user from the database
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return err
	}

	// Check if the user's email is already verified
	if user.IsVerifiedLogin {
		return c.JSON(fiber.Map{"message": "Email already verified"})
	}

	if code != user.OTPLogin {
		return c.JSON(fiber.Map{"message": "Invalid verification code"})
	}

	// Update the user's verification status in the database
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"is_verified_login": true}}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login verified successfully"})
}

func CheckUserExist(c *fiber.Ctx, email string) bool {
	if err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Err(); err != nil {
		return false
	}
	return true
}

// func Login(c *fiber.Ctx) error {
// 	// Xác thực username/password
// 	// ...

// 	// Kiểm tra và gửi OTP đến email
// 	// ...

// 	// Trả về accessToken và refreshToken
// 	// ...

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken})
// }

// func VerifyOTP(c *fiber.Ctx) error {
// 	// Xác thực OTP
// 	// ...

// 	// Trả về accessToken và refreshToken
// 	// ...

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken})
// }

// func RefreshToken(c *fiber.Ctx) error {
// 	// Xác thực refreshToken
// 	// ...

// 	// Tạo và trả về accessToken mới
// 	// ...

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": newAccessToken})
// }
