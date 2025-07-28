docker-up:
	docker-compose up --build db redis app

docker-down:
	docker-compose down

migrate:
	docker-compose run --rm migrate

migrate-test:
	docker-compose run --rm migrate-test

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

clean:
	docker-compose down -v --remove-orphans

run:
	go run cmd/main.go
