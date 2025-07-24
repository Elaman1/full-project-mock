docker-up:
	docker-compose up --build db redis app

migrate:
	docker-compose run --rm migrate

migrate-test:
	docker-compose run --rm migrate-test

docker-down:
	docker-compose down

test:
	go test ./...