// main.go
package main

import (
	"log"

	"user_login/config"
	handlers "user_login/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	//run database
	config.ConnectDB()

	// Define routes
	app.Post("/register", handlers.Register)
	app.Get("/verify", handlers.VerifyEmail)
	app.Post("/login", handlers.Login)
	app.Get("/verifyLogin", handlers.VerifyLogin)
	// app.Post("/verify", handlers.VerifyOTP)
	// app.Post("/refresh", handlers.RefreshToken)
	// app.Get("/user", handlers.GetUser)

	log.Fatal(app.Listen(":3000"))
}
