# Primo Product API

REST API service for managing products, built with Go (Fiber) + PostgreSQL + Clean Architecture.

## Repository

- **GitHub:** https://github.com/klui-work/product-service
- **Branch:** `develop`
- **Clone:**
  ```bash
  git clone https://github.com/klui-work/product-service.git
  cd product-service
  ```

## Tech Stack

- **Go 1.25** / Fiber v2
- **PostgreSQL 16** / GORM
- **Manual DI** (composition root)
- **Swagger** (gofiber/swagger)
- **Testcontainers** (integration tests)

## Project Structure

```
code/
├── main.go                          # Entry point
├── docs/                            # Swagger generated docs
├── internal/
│   ├── product/                     # Feature module
│   │   ├── domain/                  # Entity + business rules + errors
│   │   ├── application/             # DTOs (used by handler & usecase)
│   │   ├── port/
│   │   │   ├── in/                  # Input port (ProductService)
│   │   │   └── out/                 # Output port (ProductRepository)
│   │   ├── usecase/                 # Use case implementation
│   │   ├── handler/                 # HTTP layer
│   │   │   ├── v1/                  # Handler + router
│   │   │   ├── mapper/              # Domain → response mapping
│   │   │   ├── validation/          # Input validation
│   │   │   └── response/            # Standard API response wrapper
│   │   └── repository/postgresql/   # Repository impl + DB entity
│   ├── di/                          # Manual dependency injection
│   ├── infra/
│   │   ├── conf/                    # Configuration
│   │   └── database/postgresql/     # DB setup, migration, seed
│   └── utils/                       # Shared utilities
└── test/
    ├── mocks/                       # Mock repository
    ├── service/                     # Service/domain unit tests (mapper, validation)
    ├── usecase/                     # Use case unit tests (orchestrator with mock)
    ├── repo/                        # Repository integration tests (testcontainers)
    └── e2e/                         # Component E2E tests (full HTTP + DB)
```

## How to Start

### Prerequisites

- Docker & Docker Compose

### Run Service

```bash
# Start the service (dev mode with hot reload)
./run.sh
```

This will:
1. Build the Go application container
2. Start PostgreSQL 16
3. Run database migration & seed
4. Start the API server on port **3333**

### Verify

```bash
# Health check
curl http://localhost:3333/product/health
```

### API Documentation (Swagger)

Open in browser: [http://localhost:3333/api-docs/index.html](http://localhost:3333/api-docs/index.html)

## API Endpoints

### GET /product

Get all products.

```bash
curl http://localhost:3333/product
```

Response:
```json
{
  "successful": true,
  "error_code": "",
  "message": "",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567801",
      "name": "iPhone 16 Pro",
      "description": "Apple iPhone 16 Pro 256GB",
      "price": 48900,
      "sale_price": 45900
    },
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567802",
      "name": "MacBook Air M3",
      "description": "Apple MacBook Air M3 15-inch 256GB",
      "price": 44900,
      "sale_price": 42900
    },
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567803",
      "name": "iPad Pro M4",
      "description": "Apple iPad Pro M4 11-inch 256GB",
      "price": 39900,
      "sale_price": 37900
    }
  ]
}
```

### GET /product/health

Health check.

```bash
curl http://localhost:3333/product/health
```

Response:
```json
{
  "successful": true,
  "error_code": "",
  "message": "",
  "data": {
    "msg": "Product API is healthy"
  }
}
```

### POST /product

Create a new product. `description` and `sale_price` are optional/nullable.

```bash
curl -X POST http://localhost:3333/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "AirPods Pro 2",
    "description": "Apple AirPods Pro 2nd Generation",
    "price": 8990,
    "sale_price": 7990
  }'
```

Create with nullable fields omitted:

```bash
curl -X POST http://localhost:3333/product \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Minimal Product",
    "price": 100
  }'
```

Response:
```json
{
  "successful": true,
  "error_code": "",
  "message": "",
  "data": {}
}
```

Error response example:
```json
{
  "successful": false,
  "error_code": "E_VALIDATION_REQUIRED",
  "message": "name is required",
  "data": {}
}
```

### PATCH /product/{id}

Partially update a product. Only send the fields you want to change.

```bash
curl -X PATCH http://localhost:3333/product/a1b2c3d4-e5f6-7890-abcd-ef1234567801 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 16 Pro Max",
    "price": 52900
  }'
```

Response:
```json
{
  "successful": true,
  "error_code": ""
}
```

Error response example:
```json
{
  "successful": false,
  "error_code": "E_NOT_FOUND",
  "message": "product not found",
  "data": {}
}
```

## Response Format Summary

| Endpoint | Success `data` | Notes |
|----------|----------------|-------|
| GET `/product` | array of products | includes `message` |
| GET `/product/health` | `{ "msg": "..." }` | includes `message` |
| POST `/product` | `{}` | empty object, never `null` |
| PATCH `/product/{id}` | *(no `data` field)* | only `successful` + `error_code` |
| Error (all endpoints) | `{}` | includes `message` + `error_code` |

## Run Tests

```bash
cd code

# All tests
go test ./test/... -v

# Service/domain unit tests (mapper, validation)
go test ./test/service/... -v

# Use case unit tests (orchestrator with mock repo)
go test ./test/usecase/... -v

# Repository integration tests (requires Docker)
go test ./test/repo/... -v

# Component E2E tests (requires Docker)
go test ./test/e2e/... -v
```

## Verify Before Submit

```bash
cd code && go test ./test/... -v
./run.sh
curl http://localhost:3333/product/health
open http://localhost:3333/api-docs/index.html
```

## Seed Data

The service seeds 3 products on startup:

| ID | Name | Price | Sale Price |
|---|---|---|---|
| a1b2c3d4-e5f6-7890-abcd-ef1234567801 | iPhone 16 Pro | 48,900 | 45,900 |
| a1b2c3d4-e5f6-7890-abcd-ef1234567802 | MacBook Air M3 | 44,900 | 42,900 |
| a1b2c3d4-e5f6-7890-abcd-ef1234567803 | iPad Pro M4 | 39,900 | 37,900 |
