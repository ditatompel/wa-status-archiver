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
	template := "templates/layouts/app"
	if hx := c.Get("Hx-Request"); hx == "true" {
		template = ""
	}
	return c.Render("templates/dashboard", fiber.Map{
		"Title": "Dashboard",
		"Uri":   "/dashboard",
	}, template)
}

func ViewContacts(c *fiber.Ctx) error {
	co := contact.NewContactRepo(database.GetDB())

	template := "templates/layouts/app"
	if hx := c.Get("Hx-Request"); hx == "true" {
		template = ""
	}

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
		"Uri":      "/contacts",
		"Contacts": contacts,
	}, template)
}

func ViewStatusUpdates(c *fiber.Ctx) error {
	su := statusupdate.NewStatusUpdateRepo(database.GetDB())

	template := "templates/layouts/app"
	if hx := c.Get("Hx-Request"); hx == "true" {
		template = ""
	}

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
		"Uri":           "/status-updates",
		"StatusUpdates": statusUpdates,
	}, template)
}
