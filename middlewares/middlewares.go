package middlewares

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/inamandev/aspire-backend-test-2022/modules/auth"
)

var (
	token string
)

func Auth(c *fiber.Ctx) error {
	if token = c.Get("Authorization"); token == "" {
		log.Println("No token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]interface{}{
			"Message": "You are not authorised to visit this resource",
			"Success": false,
		})
	} else {
		// we have to validate the token here
		log.Println("token is here", token[7:])
		claims := auth.UserTokenClaims{}
		claims.ValidateToken(token[7:])
		log.Println(claims)
	}
	return c.Next()
}

func Recover(c *fiber.Ctx) error {
	defer func(c *fiber.Ctx) error {
		if ok := recover(); ok != nil {
			log.Println("we got a panic", ok)
			debug.PrintStack()
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
				"Message": "Something went wrong! Pleae try again later",
				"Success": false,
			})
		}
		return nil
	}(c)
	return c.Next()
}
