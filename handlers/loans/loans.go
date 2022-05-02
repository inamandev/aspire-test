package loans

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	log.Println(string(c.Body()))
	// just testing the panic
	/* if true {
		panic("just panic for testing")
	} */
	return c.SendStatus(200)
}
