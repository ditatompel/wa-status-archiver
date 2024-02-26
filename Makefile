.PHONY: static
static:
	go generate ./...

.PHONY: build
build: static
