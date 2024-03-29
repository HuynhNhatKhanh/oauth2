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

// Register a new user
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

	newUser.Password = utils.HashString(newUser.Password)
	newUser.IsVerified = false
	newUser.CreatedAt = time.Now()

	_, err := models.UserCollection.InsertOne(context.Background(), newUser)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving user"})
	}

	emailParsed := newUser.Email
	emailParsed = utils.HashString(emailParsed)

	//Send email verification
	errMail := utils.SendLinkOrOTP(newUser.Email, "link", emailParsed, newUser.Username)
	if errMail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send link"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Send link successfully"})
}

// VerifyEmail verifies the email of a user
func VerifyEmail(c *fiber.Ctx) error {
	// Get email and username from query
	emailParsed := c.Query("email")
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

	// Check if the email verification
	if !utils.CompareStringHash(emailParsed, user.Email) {
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

// Login a user
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
	if !utils.CompareStringHash(userFromDB.Password, existingUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Incorrect password"})
	}

	// Check if user is verified
	if !userFromDB.IsVerified {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email not verified"})
	}

	// Check if user is verified login
	if !userFromDB.IsVerifiedLogin {
		otp := utils.GenerateOTP()

		errMail := utils.SendLinkOrOTP(userFromDB.Email, "otp", otp, "")
		if errMail != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send OTP"})
		}
		models.UserCollection.UpdateOne(context.Background(), bson.M{"email": userFromDB.Email}, bson.M{"$set": bson.M{"otp_login": otp}})

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Send OTP to verify login"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Send OTP login failed"})
}

// VerifyLogin verifies the login of a user
func VerifyLogin(c *fiber.Ctx) error {
	// Get email and otp login from query
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

	// Check otp login
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

	// Create accessToken time live 15 minutes
	accessToken, err := utils.GenerateToken("access", time.Minute*15, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error generating access token")
	}

	// Create refreshToken time live 1 hour
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

// RefreshToken refreshes the access token
func RefreshToken(c *fiber.Ctx) error {
	// Get refresh token from request header and parse it
	refreshToken := c.Get("refreshToken")
	tokenRe, err := utils.ParseToken(refreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	// Check if the refresh token is valid
	if tokenRe.Valid {
		claims := tokenRe.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		user := models.User{}

		// Retrieve user from the database
		err := models.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		// Generate new access token time live 15 minutes
		accessToken, err := utils.GenerateToken("access", time.Minute*15, user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating access token"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Refresh token is invalid"})
	}
}

// CheckUserExist checks if a user already exists in the database
func CheckUserExist(c *fiber.Ctx, email string) bool {
	if err := models.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Err(); err != nil {
		return false
	}
	return true
}

// CheckUsernameExist checks if a username already exists in the database
func CheckUsernameExist(c *fiber.Ctx, username string) bool {
	if err := models.UserCollection.FindOne(context.Background(), bson.M{"username": username}).Err(); err != nil {
		return false
	}
	return true
}
