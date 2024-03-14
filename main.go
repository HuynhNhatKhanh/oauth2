// main.go
package main

import (
	"log"

	"user_login/config"
	handlers "user_login/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New()

	//Run database
	config.ConnectDB()

	// -------------Define routes----------------
	// Register routes
	app.Post("/register", handlers.Register)
	app.Get("/verify", handlers.VerifyEmail)
	// Login routes
	app.Post("/login", handlers.Login)
	app.Get("/verifyLogin", handlers.VerifyLogin)
	// User routes
	app.Get("/user", handlers.GetUser)
	// Refresh token
	app.Post("/refresh", handlers.RefreshToken)

	log.Fatal(app.Listen(":3000"))
}
