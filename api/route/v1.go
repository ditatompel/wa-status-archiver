package route

import (
	"wabot/api/controller"

	"github.com/gofiber/fiber/v2"
)

func V1Api(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Post("/auth/login", controller.Login)
	v1.Delete("/auth/login", controller.Logout)
}
