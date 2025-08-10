up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build

logs:
	docker-compose logs -f

migrate-up:
	docker run --rm -v $(PWD)/deploy/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://logscale:logscale123@postgres:5432/logscale?sslmode=disable" up

migrate-down:
	docker run --rm -v $(PWD)/deploy/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://logscale:logscale123@postgres:5432/logscale?sslmode=disable" down 1

# Alternative: Run migrations directly in the postgres container
migrate-sql:
	docker-compose exec postgres psql -U logscale -d logscale -f /docker-entrypoint-initdb.d/001_create_logs_table.sql

# Development commands
dev:
	go run cmd/api/main.go

.PHONY: test test-coverage test-integration

test:
	go test ./test -v

test-coverage:
	go test ./test -v -cover

test-integration:
	go test ./test -v -tags=integration

clean:
	docker-compose down -v
	docker system prune -f
