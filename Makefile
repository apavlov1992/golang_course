build: ## Build binary file
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o myapp myapp.go