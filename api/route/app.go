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
	app.Get("/contacts/hxp", middleware.CookieProtected, controller.ViewContactPartials)
	app.Get("/status-updates", middleware.CookieProtected, controller.ViewStatusUpdates)
	app.Get("/status-updates/hxp", middleware.CookieProtected, controller.ViewStatusUpdatePartials)
}
