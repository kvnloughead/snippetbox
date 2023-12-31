run-dev:
	export GO_ENV=development
	go run ./cmd/web

run-prod:
	export GO_ENV=production
	go run ./cmd/web

build:
	go build -o snippetbox-compiled

test:
	go test ./...

.PHONY: run-dev run-prod build test