up:
	docker-compose up -d

down:
	docker-compose down

migrate-up:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://admin:secret@db:5432/logsdb?sslmode=disable" up

migrate-down:
	docker run --rm -v $(PWD)/db/migrations:/migrations --network logscale_default migrate/migrate -path=/migrations -database "postgres://admin:secret@db:5432/logsdb?sslmode=disable" down 1
