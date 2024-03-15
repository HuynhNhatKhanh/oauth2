// handlers/auth_handler.go
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"user_login/models"
	"user_login/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

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

	// Check if username already exists
	if CheckUsernameExist(c, newUser.Username) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Username already exists"})
	}

	// otp := utils.GenerateOTP()
	// newUser.OTP = otp
	newUser.Password = utils.HashPassword(newUser.Password)
	newUser.IsVerified = false
	newUser.CreatedAt = time.Now()

	_, err := models.UserCollection.InsertOne(context.Background(), newUser)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving user"})
	}

	emailParsed := newUser.Email
	emailParsed = utils.HashPassword(emailParsed)

	//Send email verification
	errMail := utils.SendOTP(newUser.Email, "link", emailParsed, newUser.Username)
	if errMail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Send otp successfully"})
}

func VerifyEmail(c *fiber.Ctx) error {
	emailParsed := c.Query("email")
	// code := c.Query("otp")
	username := c.Query("username")

	fmt.Println(emailParsed)
	fmt.Println(username)

	// Retrieve user from the database
	var user models.User
	err := models.UserCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return err
	}

	// Check if the user's email is already verified
	if user.IsVerified {
		return c.JSON(fiber.Map{"message": "Email already verified"})
	}

	// if code != user.OTP {
	// 	return c.JSON(fiber.Map{"message": "Invalid verification code"})
	// }

	if !utils.ComparePasswordHash(emailParsed, user.Email) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email verified failed"})
	}

	// Update the user's verification status in the database
	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": bson.M{"is_verified": true}}
	_, err = models.UserCollection.UpdateOne(context.Background(), filter, update)
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
	result := models.UserCollection.FindOne(context.Background(), filter)
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

		errMail := utils.SendOTP(userFromDB.Email, "otp", otp, "")
		if errMail != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP"})
		}
		models.UserCollection.UpdateOne(context.Background(), bson.M{"email": userFromDB.Email}, bson.M{"$set": bson.M{"otp_login": otp}})

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Send OTP to verify login"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Send OTP login failed"})
}

func VerifyLogin(c *fiber.Ctx) error {
	email := c.Query("email")
	code := c.Query("otp_login")

	// Retrieve user from the database
	var user models.User
	err := models.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
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
	_, err = models.UserCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	// Create accessToken and refreshToken
	accessToken, err := utils.GenerateToken("access", time.Minute*15, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error generating access token")
	}

	refreshToken, err := utils.GenerateToken("refresh", time.Hour, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error generating refresh token")
	}
	response := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func CheckUserExist(c *fiber.Ctx, email string) bool {
	if err := models.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Err(); err != nil {
		return false
	}
	return true
}

func CheckUsernameExist(c *fiber.Ctx, username string) bool {
	if err := models.UserCollection.FindOne(context.Background(), bson.M{"username": username}).Err(); err != nil {
		return false
	}
	return true
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("refreshToken")
	tokenRe, err := utils.ParseToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}
	if tokenRe.Valid {
		claims := tokenRe.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		user := models.User{}
		err := models.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		accessToken, err := utils.GenerateToken("access", time.Minute*15, user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating access token"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Refresh token is invalid"})
	}
}
