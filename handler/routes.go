package handler

import (
	"github.com/gofiber/fiber/v2"
)

func AppRoute(app *fiber.App) {
	app.Get("/", LoginView)
	app.Post("/", Login)
	app.Delete("/", Logout)
	app.Get("/dashboard", CookieProtected, DashboardView)
	app.Get("/contacts", CookieProtected, ContactsView)
	app.Get("/contacts/hxp", CookieProtected, ContactPartials)
	app.Get("/status-updates", CookieProtected, StatusUpdatesView)
	app.Get("/status-updates/hxp", CookieProtected, StatusUpdatePartials)
}
