package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:generate npm i
//go:generate npm run build
//go:embed all:build/*
var f embed.FS

func SvelteKitHandler() http.FileSystem {
	build, err := fs.Sub(f, "build")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(build)
}
