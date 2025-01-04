include .env

all: build test

build:
	@echo "Building..."
	@go build -o main cmd/http/main.go

clean:
	@echo "Cleaning..."
	@rm -f main

# Run integrations tests
itest:
	@echo "Running integration tests..."
	@go test --tags=integration -v ./...

migration-down:
	@cd internal/adapters/storage/postgres/migrations && goose postgres "$(DB_ADDR)" down

migration-up:
	@cd internal/adapters/storage/postgres/migrations && goose postgres "$(DB_ADDR)" up

test:
	@echo "Testing..."
	@go test ./... -v

run:
	@go run cmd/http/main.go

seed:
	@go run cmd/seed/main.go

swag:
	@swag init -g cmd/http/main.go -o ./docs

# Live reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY:  all build clean itest migration-down migration-reset migration-up test run seed swag watch