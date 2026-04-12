dev:
	docker compose up --build

down:
	docker compose down

clear-down:
	docker compose down -v

test:
	go test ./...

lint:
	golangci-lint run ./...