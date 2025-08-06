docker-up:
	docker-compose up --build db redis prometheus grafana app

docker-down:
	docker-compose down -v

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
