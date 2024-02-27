package controller

import (
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

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Logged in",
		"data":    nil,
	})
}
