.PHONY: run test migrate tidy docker-up docker-down

run:
	go run cmd/gateway/main.go

test:
	go test -v ./internal/router/...

tidy:
	go mod tidy

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
