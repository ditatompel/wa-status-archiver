package api

import (
	"encoding/json"
	"fmt"
	"time"
	"wabot/internal/database"
	"wabot/internal/repo"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	payload := repo.Admin{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	admin := repo.NewAdminRepo(database.GetDB())
	res, err := admin.Login(payload.Username, payload.Password)
	if err != nil {
		triggerJson, _ := json.Marshal(map[string]interface{}{"err": err.Error()})

		c.Response().Header.Set("HX-Trigger", string(triggerJson))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	token := fmt.Sprintf("auth_%d_%d", res.Id, time.Now().Unix())
	c.Cookie(&fiber.Cookie{
		Name:     "wabot_ui",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	c.Response().Header.Set("HX-Redirect", "/dashboard")

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Logged in",
		"data":    nil,
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "wabot_ui",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
	})

	c.Response().Header.Set("HX-Redirect", "/")

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Logged out",
		"data":    nil,
	})
}

func ViewLogin(c *fiber.Ctx) error {
	cookie := c.Cookies("wabot_ui")
	if cookie != "" {
		return c.Redirect("/dashboard")
	}
	return c.Render("templates/index", fiber.Map{
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
	co := repo.NewContactRepo(database.GetDB())

	template := "templates/layouts/app"
	if hx := c.Get("Hx-Request"); hx == "true" {
		template = ""
	}

	query := repo.ContactQueryParams{
		Search:      c.Query("search"),
		Sort:        c.Query("sort"),
		Dir:         c.Query("dir"),
		RowsPerPage: 20,
		Page:        c.QueryInt("page", 1),
	}

	contacts, err := co.Contacts(query)
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
		"NextPage": contacts.NextPage,
		"Contacts": contacts.Contacts,
	}, template)
}

func ViewContactPartials(c *fiber.Ctx) error {
	co := repo.NewContactRepo(database.GetDB())

	query := repo.ContactQueryParams{
		Search:      c.Query("search"),
		Sort:        c.Query("sort"),
		Dir:         c.Query("dir"),
		RowsPerPage: 20,
		Page:        c.QueryInt("page", 1),
	}

	contacts, err := co.Contacts(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Render("templates/partials/contact", fiber.Map{
		"Search":   query.Search,
		"NextPage": contacts.NextPage,
		"Contacts": contacts.Contacts,
	})
}

func ViewStatusUpdates(c *fiber.Ctx) error {
	su := repo.NewStatusUpdateRepo(database.GetDB())
	query := repo.StatusUpdateQueryParams{
		JID:         c.Query("jid"),
		Sort:        c.Query("sort"),
		Dir:         c.Query("dir"),
		RowsPerPage: 10,
		Page:        c.QueryInt("page", 1),
	}

	template := "templates/layouts/app"
	if hx := c.Get("Hx-Request"); hx == "true" {
		template = ""
	}

	statusUpdates, err := su.StatusUpdates(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}
	contacts, _ := su.Contacts()

	return c.Render("templates/status-updates", fiber.Map{
		"Title":         "Status Updates",
		"Uri":           "/status-updates",
		"NextPage":      statusUpdates.NextPage,
		"Contacts":      contacts,
		"JID":           query.JID,
		"StatusUpdates": statusUpdates.Statuses,
	}, template)
}

func ViewStatusUpdatePartials(c *fiber.Ctx) error {
	su := repo.NewStatusUpdateRepo(database.GetDB())
	query := repo.StatusUpdateQueryParams{
		JID:         c.Query("jid"),
		Sort:        c.Query("sort"),
		Dir:         c.Query("dir"),
		RowsPerPage: 10,
		Page:        c.QueryInt("page", 1),
	}

	statusUpdates, err := su.StatusUpdates(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Render("templates/partials/statuses", fiber.Map{
		"JID":           query.JID,
		"NextPage":      statusUpdates.NextPage,
		"StatusUpdates": statusUpdates.Statuses,
	})
}
