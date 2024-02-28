package views

import (
	"embed"
	"net/http"
)

//go:embed all:templates/*
var f embed.FS

func EmbedHandler() http.FileSystem {
	return http.FS(f)
}
