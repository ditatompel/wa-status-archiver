BINARY_NAME = wabot

.PHONY: build
build: ui linux64

.PHONY: ui
ui:
	go generate ./...

.PHONY: linux64
linux64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-amd64

.PHONY: clean
clean:
	go clean
	rm -rfv ./bin
