package controller

import (
	"wabot/internal/database"
	"wabot/internal/repo/contact"
	"wabot/internal/repo/statusupdate"

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
	co := contact.NewContactRepo(database.GetDB())

	contacts, err := co.Contacts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Render("templates/contacts", fiber.Map{
		"Title":    "Contacts",
		"Uri":      "contacts",
		"Contacts": contacts,
	}, "templates/layouts/app")
}

func ViewStatusUpdates(c *fiber.Ctx) error {
	su := statusupdate.NewStatusUpdateRepo(database.GetDB())

	statusUpdates, err := su.StatusUpdates()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Render("templates/status-updates", fiber.Map{
		"Title":         "Status Updates",
		"Uri":           "status-updates",
		"StatusUpdates": statusUpdates,
	}, "templates/layouts/app")
}
