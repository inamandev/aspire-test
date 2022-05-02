package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/golang-jwt/jwt/v4"
	_ "github.com/google/uuid"
	"github.com/inamandev/aspire-backend-test-2022/handlers/auth"
	"github.com/inamandev/aspire-backend-test-2022/handlers/loans"
	"github.com/inamandev/aspire-backend-test-2022/middlewares"
	_ "github.com/mattn/go-sqlite3"
)

var (
	app *fiber.App
)

/* func init() {
	if ok := godotenv.Load(); ok != nil {
		log.Println(ok)
		panic(ok)
	}
} */

// import _
func main() {
	defer handlePanic()
	app = fiber.New()
	app.Use(etag.New())
	app.Use(logger.New())
	app.Use(middlewares.Recover)
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/user/auth", auth.Login)
	v1.Post("/loans", middlewares.Auth, loans.Create)
	if ok := app.Listen(":3031"); ok != nil {
		log.Println("Unable to start the application due to below error")
		log.Println(ok)
		panic(ok)
	}
}

func handlePanic() {
	if ok := recover(); ok != nil {
		log.Println("recovering from panic", ok)
		app.Listen(":3031")
	}
}
