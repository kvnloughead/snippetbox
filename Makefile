run-dev:
	export GO_ENV=development && export PORT=4001 && air

run-prod:
	export GO_ENV=production
	go run ./cmd/web

build:
	go build -o ./tmp/myserver ./cmd/web

test:
	go test ./...

.PHONY: run-dev run-prod build test