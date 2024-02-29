package api

import (
	"github.com/gofiber/fiber/v2"
)

func CookieProtected(c *fiber.Ctx) error {
	cookie := c.Cookies("wabot_ui")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	return c.Next()
}
