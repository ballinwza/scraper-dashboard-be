# Go (Gin) Backend Service

[https://scraper-dashboard-be-prod-758337397665.asia-southeast1.run.app](https://scraper-dashboard-be-prod-758337397665.asia-southeast1.run.app)
[Frontend Service](https://scraper-dashboard-fe-prod-758337397665.asia-southeast1.run.app)
## Code Structure & Clean Architecture
```
github.com/ballinwza/scraper-dashboard-be/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point: Dependency Injection & Start HTTP Server
│
├── config/                          # Configuration loader
│
├── internal/                        # Private application and code
│   ├── domain/                      # Domain Layer (Pure Go Entities & Interfaces)
│   │
│   ├── usecase/                     # Usecase Layer (Business Logic)
│   │
│   ├── delivery/                    # 3. Delivery Layer (Adapters for Inbound Calls)
│   │   ├── http/                    # REST API Handlers (Gin / Fiber)
│   │   │   ├── handler/
│   │   │   ├── middleware/          # Auth JWT, CORS, Logging
│   │   │   └── router.go            # Route registrations
│   │   ├── grpc/         
│   │   │   ├── handler/             # Implements gRPC Generated Interfaces
│   │   │   ├── interceptor/         # Middleware สำหรับ gRPC (Auth, Logging, Recovery)
│   │   │   └── server.go            # gRPC Server Initialization & Registration
│   │   └── ws/                      # WebSockets for Real-time Scraper Logs
│   │
│   └── repository/                  # Infrastructure Layer (Adapters for Outbound Calls)
│       ├── mongodb/                
│       ├── scraper/                
│       └── redis/                   
│
├── pkg/                             # Public libraries reusable by other microservices  
│
└── docs/                            # Swagger / OpenAPI documentation
```

## Dependencies
* Gin - go get github.com/gin-gonic/gin
* MongoDB - go get go.mongodb.org/mongo-driver/v2/mongo
* Gorilla - websocket go get github.com/gorilla/websocket
* Viper - go get github.com/spf13/viper
* godotenv - go get github.com/joho/godotenv
* JWT - go get github.com/golang-jwt/jwt/v5
* Swagger - go install github.com/swaggo/swag/cmd/swag@latest
* Gin swaggo - go get github.com/swaggo/gin-swagger
* Swaggo file - go get github.com/swaggo/files
* Zap Logger - go get go.uber.org/zap
* Validator - go get github.com/go-playground/validator/v10
* Colly scraper - github.com/gocolly/colly/v2