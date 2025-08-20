up:
	docker-compose up -d

down:
	docker-compose down
	docker stop $$(docker ps -q --filter "name=logscale") 2>/dev/null || true

build:
	docker-compose build

logs:
	docker-compose logs -f

# Database migrations
migrate-up:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://logscale:logscale123@postgres:5432/logscale?sslmode=disable" up

migrate-down:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://logscale:logscale123@postgres:5432/logscale?sslmode=disable" down 1

migrate-force:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://logscale:logscale123@postgres:5432/logscale?sslmode=disable" force $(version)

migrate-version:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://logscale:logscale123@postgres:5432/logscale?sslmode=disable" version

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

# Development commands
dev:
	go run cmd/api/main.go

dev-worker:
	go run cmd/worker/main.go

.PHONY: test test-coverage test-integration lint

lint:
	golangci-lint run --no-config

test:
	go test ./test -v

test-integration:
	go test ./test -v -tags=integration

test-coverage:
	go test ./test -v -cover

clean:
	docker-compose down -v
	docker stop $$(docker ps -q --filter "name=logscale") 2>/dev/null || true
	docker rm $$(docker ps -aq --filter "name=logscale") 2>/dev/null || true
	docker system prune -f
