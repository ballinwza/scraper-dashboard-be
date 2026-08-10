swag:
	swag init -g cmd/server/main.go

dev:
	go run cmd/server/main.go

clean:
	rm -rf tmp
	go clean

tidy:
	go mod tidy

air-init:
	air init

air-dev:
	air

docker-build:
	docker build -t scraper-dashboard-be .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env scraper-dashboard-be

docker-clear:
	docker run --rm -it scraper-dashboard-be /bin/sh

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down -v