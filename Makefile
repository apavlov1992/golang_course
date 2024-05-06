# Build targets
.PHONY: build/xkcd build/stemming

build: ## Build binary file for xkcd
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o xkcd internal/xkcd/xkcd.go

build/stemming: ## Build binary file for stemming
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o stemming internal/stemming/stemming.go

build: ## Build main projects
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main cmd/main.go

