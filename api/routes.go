package api

import (
	"github.com/gofiber/fiber/v2"
)

func AppRoute(app *fiber.App) {
	app.Get("/", ViewLogin)
	app.Get("/dashboard", CookieProtected, ViewDashboard)
	app.Get("/contacts", CookieProtected, ViewContacts)
	app.Get("/contacts/hxp", CookieProtected, ViewContactPartials)
	app.Get("/status-updates", CookieProtected, ViewStatusUpdates)
	app.Get("/status-updates/hxp", CookieProtected, ViewStatusUpdatePartials)
}

func V1Api(app *fiber.App) {
	v1 := app.Group("/api/v1")
	v1.Post("/auth/login", Login)
	v1.Delete("/auth/login", Logout)
}
