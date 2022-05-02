package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/inamandev/aspire-backend-test-2022/modules/auth"
	u "github.com/inamandev/aspire-backend-test-2022/modules/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	user      u.User
	err       error
	validAuth error
)

func Login(c *fiber.Ctx) error {
	var token string
	userLoginDetails := auth.UserAuth{}
	if ok := c.BodyParser(&userLoginDetails); ok != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{
			"success": false,
			"message": "your parameters are not correct",
		})
	}
	if user, err = u.GetByUsername(userLoginDetails.Username); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]interface{}{
			"success": false,
			"message": "username or password is incorrect",
		})
	}
	if validAuth = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLoginDetails.Password)); validAuth != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]interface{}{
			"success": false,
			"message": "username or password is incorrect",
		})
	}
	if validAuth == nil {
		// User authentication successful. create Token
		claims := auth.UserTokenClaims{
			Id:     user.Id,
			Role:   user.Role,
			Status: user.Status,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: "Naman Kathuria",
				ExpiresAt: &jwt.NumericDate{
					Time: time.Now().Add(time.Hour * 4),
				},
			},
		}
		if token, err = claims.CreateToken(); err != nil {
			log.Println("err during creating jwt token", err)
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
				"success": false,
				"message": "Someting went wrong! Please try again",
			})
		}
		if err == nil {
			c.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
				"success": true,
			})
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
		"success": false,
		"message": "Someting went wrong! Please try again",
	})
}
