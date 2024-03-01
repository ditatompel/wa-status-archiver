BINARY_NAME = wabot

.PHONY: tailwind
tailwind:
	npx tailwindcss -i ./views/css/main.css -o ./public/main.css --minify

.PHONY: build
build: tailwind linux64

.PHONY: linux64
linux64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-amd64

.PHONY: clean
clean:
	go clean
	rm -rfv ./bin
