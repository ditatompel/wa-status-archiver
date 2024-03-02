.PHONY: tailwind copyhtmx build linux64

BINARY_NAME = wa-status-archiver

build: copyhtmx tailwind linux64

tailwind:
	npx tailwindcss -i ./views/css/main.css -o ./public/main.css --minify

copyhtmx:
	cp ./node_modules/htmx.org/dist/htmx.min.js ./public

linux64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-amd64

clean:
	go clean
	rm -rfv ./bin
	rm ./public/htmx.min.js
	rm ./public/main.css
