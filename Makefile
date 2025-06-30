# Migrations 
migrate:
	go run cmd/migrate/main.go

# Rollback migrations
rollback:
	go run cmd/migrate/main.go rollback

# Run the application
run:
	go run cmd/server/main.go

# Generate docs
.PHONY: docs
docs:
	swag init -g cmd/server/main.go

# Run tests
.PHONY: test
test:
	go test -race -v -cover  .\...

# Run worker
.PHONY: worker
worker:
	go run cmd/worker/main.go


# Run dev docker image 
.PHONY: dev
dev:
	docker-compose -f docker\compose.dev.yml up