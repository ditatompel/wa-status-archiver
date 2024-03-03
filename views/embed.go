package views

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed all:templates/*
var templates embed.FS

func EmbedTemplates() http.FileSystem {
	return http.FS(templates)
}

//go:embed static/*
var embedStatic embed.FS

func EmbedStatic() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root:       http.FS(embedStatic),
		PathPrefix: "static",
		Browse:     true,
	})
}
