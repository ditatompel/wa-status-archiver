package controller

import (
	"github.com/gofiber/fiber/v2"
)

func ViewLogin(c *fiber.Ctx) error {
	cookie := c.Cookies("wabot_ui")
	if cookie != "" {
		return c.Redirect("/dashboard")
	}
	return c.Render("templates/login", fiber.Map{
		"Title": "Login",
		"Uri":   "login",
	}, "templates/layouts/main")
}

func ViewDashboard(c *fiber.Ctx) error {
	return c.Render("templates/dashboard", fiber.Map{
		"Title": "Dashboard",
		"Uri":   "dashboard",
	}, "templates/layouts/app")
}

func ViewContacts(c *fiber.Ctx) error {
	return c.Render("templates/contacts", fiber.Map{
		"Title": "Contacts",
		"Uri":   "contacts",
	}, "templates/layouts/app")
}
