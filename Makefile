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

.PHONY: proto clean

# กำหนดโฟลเดอร์ปลายทาง
PROTO_DIR = internal/delivery/proto
PROTO_OUT_DIR = internal/delivery/grpc/api


# Target สำหรับ Generate Protobuf ทั้งหมดในโฟลเดอร์ api/proto/v1
proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "✅ Generated gRPC code successfully!"

# Target สำหรับลบไฟล์ที่ generate ออกมา
clean-proto:
	go run -v github.com/google/go-licenses@latest --help >nul 2>&1 || true
	powershell -Command "Remove-Item -Path '$(PROTO_OUT_DIR)/*.pb.go' -ErrorAction SilentlyContinue"
	@echo "🧹 Cleaned generated proto files."

rebuild-proto:
	clean-proto proto
