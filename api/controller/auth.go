package controller

import (
	"encoding/json"
	"fmt"
	"time"
	"wabot/internal/database"
	"wabot/internal/repo/admin"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	payload := admin.Admin{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	repo := admin.NewAdminRepo(database.GetDB())
	res, err := repo.Login(payload.Username, payload.Password)
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

	c.Response().Header.Set("HX-Redirect", "/login")

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Logged out",
		"data":    nil,
	})
}
