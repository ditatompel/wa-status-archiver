package route

import (
	"wabot/api/controller"
	"wabot/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppRoute(app *fiber.App) {
	app.Get("/login", controller.ViewLogin)
	app.Get("/dashboard", middleware.CookieProtected, controller.ViewDashboard)
	app.Get("/contacts", middleware.CookieProtected, controller.ViewContacts)
	app.Get("/status-updates", middleware.CookieProtected, controller.ViewStatusUpdates)
}
