BINARY_NAME = wabot

.PHONY: build
build: linux64

.PHONY: linux64
linux64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-static-linux-amd64

.PHONY: clean
clean:
	go clean
	rm -rfv ./bin
