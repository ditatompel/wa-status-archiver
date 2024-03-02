package handler

import (
	"github.com/gofiber/fiber/v2"
)

func AppRoute(app *fiber.App) {
	app.Get("/", ViewLogin)
	app.Post("/", Login)
	app.Delete("/", Logout)
	app.Get("/dashboard", CookieProtected, ViewDashboard)
	app.Get("/contacts", CookieProtected, ViewContacts)
	app.Get("/contacts/hxp", CookieProtected, ViewContactPartials)
	app.Get("/status-updates", CookieProtected, ViewStatusUpdates)
	app.Get("/status-updates/hxp", CookieProtected, ViewStatusUpdatePartials)
}
