package handler

import (
	"github.com/gofiber/fiber/v2"
)

func CookieProtected(c *fiber.Ctx) error {
	cookie := c.Cookies("wabot_ui")
	if cookie == "" {
		c.Response().Header.Set("HX-Redirect", "/")
		return c.Redirect("/")
	}

	return c.Next()
}
