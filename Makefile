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

