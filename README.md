# Cinema Booking System

A **concurrent, production-ready cinema seat booking platform** built with Go microservices. The system leverages **gRPC** for inter-service communication, **PostgreSQL** for persistent storage, **Redis** for distributed caching and seat-lock management, and a **Gin-based API Gateway** as the single entry point for all client requests.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                  Client (Browser/App)               │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP REST
┌──────────────────────▼──────────────────────────────┐
│               API Gateway  (:3333)                  │
│  Gin · CORS · Security Headers · Request ID         │
│  Logging · Error Handling · Prometheus Metrics      │
└────┬──────────────────┬───────────────────┬─────────┘
     │ gRPC             │ gRPC              │ gRPC
┌────▼────┐       ┌─────▼─────┐      ┌─────▼─────┐
│  Users  │       │ Bookings  │      │   Seats   │
│ Service │       │  Service  │      │  Service  │
│ :50051  │       │  :50052   │      │  :50053   │
└────┬────┘       └─────┬─────┘      └─────┬─────┘
     │                 │                   │
┌────▼────┐       ┌─────▼─────┐      ┌─────▼─────┐
│Postgres │       │ Postgres  │      │ Postgres  │
│user_svc │       │booking_svc│      │ seats_svc │
└─────────┘       └─────┬─────┘      └─────────-─┘
                        │
                  ┌─────▼─────┐
                  │   Redis   │
                  │(seat lock)│
                  └───────────┘
```

---

## Services

| Service | Port | Protocol | Responsibility |
|---|---|---|---|
| **Gateway** | `3333` | HTTP REST | Routes client requests to downstream gRPC services |
| **Users** | `50051` | gRPC | User registration and authentication |
| **Bookings** | `50052` | gRPC | Booking lifecycle (create, update, cancel, query) |
| **Seats** | `50053` | gRPC | Seat inventory management |

All services share a **`common`** module that holds the Protobuf-generated Go code (`api/`) for type-safe gRPC contracts.

---

## Project Structure

```
cinema_booking/
├── go.work                  # Go workspace — ties all modules together
├── Makefile                 # Top-level proto generation
│
├── common/                  # Shared Protobuf generated code
│   └── api/                 # booking.pb.go, booking_grpc.pb.go, ...
│
├── gateway/                 # API Gateway (Gin HTTP server)
│   ├── main.go
│   ├── http_handler.go      # Route registration & HTTP handlers
│   ├── config/
│   └── middlewares/         # Auth, RequestID, Logging, Error handling
│
├── users/                   # Users microservice
│   ├── cmd/
│   │   ├── main.go
│   │   └── migrate/         # DB migration entry point
│   ├── internal/
│   │   ├── handler/         # gRPC handler
│   │   ├── model/           # User struct
│   │   ├── repository/      # GORM data access layer
│   │   └── service/         # Business logic (hashing, email check)
│   └── pkg/
│       ├── helpers/         # Password hashing utilities
│       └── type/            # Service & repository interfaces
│
├── bookings/                # Bookings microservice
│   ├── cmd/
│   │   ├── main.go
│   │   └── migrate/
│   ├── internal/
│   │   ├── handler/         # gRPC handler (Create, FindAll, IsSeatAvailable, Update, Cancel)
│   │   ├── model/           # Booking struct
│   │   ├── repository/      # GORM data access layer
│   │   └── service/         # Business logic & seat-availability check
│   └── pkg/types/           # Service & repository interfaces
│
└── seats/                   # Seats microservice
    ├── cmd/
    │   ├── main.go
    │   └── migrate/
    ├── internal/
    │   ├── handler/         # gRPC handler (SetSeats, GetSeats, DeleteSeat)
    │   ├── model/           # Seat struct (ID, SeatNumber, Row, Column)
    │   ├── repository/      # GORM data access layer
    │   └── service/         # Seat seeding and management
    └── pkg/types/           # Service & repository interfaces
```

---

## Technology Stack

| Category | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| RPC Framework | [gRPC](https://grpc.io/) + Protocol Buffers |
| ORM | [GORM](https://gorm.io/) |
| Database | PostgreSQL |
| Cache / Locks | Redis |
| Auth | JWT |
| Observability | [Prometheus](https://prometheus.io/) + [Logrus](https://github.com/sirupsen/logrus) |
| Hot Reload | [Air](https://github.com/air-verse/air) |
| Security | `unrolled/secure` (XSS, CSP, frame-deny) |

---

## API Endpoints (Gateway)

### Users

| Method | Path | Description |
|---|---|---|
| `POST` | `/users` | Register a new user |
| `POST` | `/users/login` | User login |

### Bookings

| Method | Path | Description |
|---|---|---|
| `POST` | `/bookings` | Create a new booking |
| `POST` | `/bookings/hold` | Hold a booking |
| `GET` | `/bookings` | Get user's bookings |
| `DELETE` | `/bookings` | Cancel booking |

### Seats

| Method | Path | Description |
|---|---|---|
| `POST` | `/seats/set` | Seating seats with mock data |
| `GET` | `/seats` | Get seats |

### System

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

---

## Configuration

Each service is configured via environment variables (loaded from `.env` in non-production environments).

### Gateway (`.env`)

```env
PORT=3333
GIN_MODE=debug                         # debug | release
ENVIRONMENT=development                # development | production
DOMAIN=localhost

USERS_SERVICE_ADDRESS=localhost:50051
BOOKINGS_SERVICE_ADDRESS=localhost:50052
SEATS_SERVICE_ADDRESS=localhost:50053
```

### Users / Bookings / Seats (`.env`)

```env
BOOKINGS_SERVICE_ADDRESS=localhost:50052
SEATS_SERVICE_ADDRESS=localhost:50053

ENVIRONMENT=development
DATABASE_HOST=localhost
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=secret
DATABASE_NAME=<service_db_name>
DATABASE_PORT=5432

JWT_SECRET_KEY=your_jwt_secret
```

> **Note:** In `production` mode, environment variables are expected to be injected by the host (e.g., Docker/Kubernetes) — `.env` files are not loaded.

---

## Data Models

### User
| Field | Type |
|---|---|
| `name` | `string` |
| `email` | `string` |
| `password` | `string` (bcrypt hashed) |

### Booking
| Field | Type |
|---|---|
| `id` | `int64` |
| `seat_id` | `int64` |
| `seat_number` | `string` |
| `user_id` | `int64` |
| `user_name` | `string` |
| `movie_id` | `int64` |
| `show_time` | `string` |

### Seat
| Field | Type |
|---|---|
| `id` | `int64` |
| `seat_number` | `string` (e.g. `A1`) |
| `row_number` | `string` (e.g. `A`) |
| `column_number` | `string` (e.g. `1`) |
| `status` | `string` (e.g. `Available`) |
| `created_at` | `time.Time` |
| `updated_at` | `time.Time` |

---

## gRPC Service Contracts

### UserService
```protobuf
service UserService {
  rpc CreateUser (CreateUserRequest) returns (CreateUserResponse);
  rpc LoginUser  (LoginUserRequest)  returns (LoginUserResponse);
}
```

### BookingService
```protobuf
service BookingService {
  rpc Create          (CreateRequest)          returns (CreateResponse);
  rpc FindAll         (FindAllRequest)         returns (FindAllResponse);
  rpc IsSeatAvailable (IsSeatAvailableRequest) returns (IsSeatAvailableResponse);
  rpc Update          (UpdateRequest)          returns (UpdateResponse);
  rpc Cancel          (CancelRequest)          returns (CancelResponse);
}
```

### SeatsService
```protobuf
service SeatsService {
  rpc SetSeats   (SetSeatsRequest)   returns (SetSeatsResponse);
  rpc GetSeats   (GetSeatsRequest)   returns (GetSeatsResponse);
  rpc DeleteSeat (DeleteSeatRequest) returns (DeleteSeatResponse);
}
```

---

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc`
- [Air](https://github.com/air-verse/air) *(optional, for hot reload)*

### 1. Clone the repository

```bash
git clone https://github.com/ThuraMinThein/cinema_booking_system.git
cd cinema_booking_system
```

### 2. Generate Protobuf code

```bash
make generate-proto
```

### 3. Set up environment variables

Copy and configure the `.env` file in each service directory:

```bash
cp gateway/.env.example gateway/.env
cp users/.env.example users/.env
cp bookings/.env.example bookings/.env
cp seats/.env.example seats/.env
```

### 4. Run database migrations

```bash
# In each service directory
go run cmd/migrate/main.go
```

### 5. Start all services

Open a terminal for each service and run:

```bash
# Gateway
cd gateway && air        # or: go run .

# Users service
cd users && air          # or: go run cmd/main.go

# Bookings service
cd bookings && air       # or: go run cmd/main.go

# Seats service
cd seats && air          # or: go run cmd/main.go
```

### 6. Verify the gateway is running

```bash
curl http://localhost:3333/health
# {"status":"ok"}
```

---

## Middleware & Security

The API Gateway applies the following middleware stack to every request:

| Middleware | Purpose |
|---|---|
| `gin.Recovery()` | Panic recovery |
| `RequestIDMiddleware` | Attaches a `X-Request-ID` UUID to every request/response |
| `LoggingMiddleware` | Structured request logging via Logrus |
| `ErrorHandlingMiddleware` | Centralized error response formatting |
| `secure` (unrolled) | Sets `X-Frame-Options`, `X-Content-Type-Options`, `X-XSS-Protection`, and `Content-Security-Policy` headers |
| CORS | Configurable origin whitelist with credential support |

---

## Observability

- **Logging:** Structured logs via [Logrus](https://github.com/sirupsen/logrus). JSON format in `release` mode, human-readable text in development.
- **Metrics:** Prometheus metrics exposed at `GET /metrics` on the gateway.
- **Tracing:** Request IDs (`X-Request-ID`) are propagated across the request lifecycle and can be correlated in logs.

---

## Concurrency & Seat Locking

The system is designed for **high-concurrency booking scenarios**:

- **`IsSeatAvailable` RPC** is called before any booking is persisted to prevent double-booking.
- **Redis** is used as a distributed lock store so that seat availability checks and reservations are atomic across multiple gateway instances.
- The `BookingService` enforces the invariant: a seat can only have one active booking per movie + show time.
---

---

## Development

```bash
# Regenerate Protobuf stubs (from workspace root)
make generate-proto

# Hot-reload a service with Air
cd <service> && air

# Run a service manually
cd <service> && go run cmd/main.go
```

---

## Performance

### Optimizations

#### Get All Seats — `GET /seats` *(14 May 2026)*

**Improvement: ~20 s → ~1.5 s average response time**

| Technique | Details |
|---|---|
| **Batch operations** | Replaced row-by-row DB & Redis lookups with bulk batch fetches to minimize round-trips |
| **Database indexes** | Added indexes on frequently queried seat columns to speed up query execution |
| **Redis key scanning** | Switched to efficient Redis key-scan patterns instead of individual key lookups per seat |

**Load-test results** (1 000 concurrent requests, concurrency 50):

| Metric | Value |
|---|---|
| Total Requests | 1 000 |
| Concurrency | 50 |
| Successes | 998 |
| Failures | 2 |
| Total Time | 24.23 s |
| Average Latency | **1 119.86 ms** |
| Requests / sec | 41.28 |

---

#### Hold Booking — `POST /bookings/hold` *(15 May 2026)*

**Improvement: 5 - 8 s → 1.5 - 2 s average response time**

| Technique | Details |
|---|---|
| **Parallel gRPC requests** | Executed user-service and seat-service gRPC calls concurrently using goroutines and sync.WaitGroup to reduce total request latency |
| **Redis distributed locking** | Replaced separate availability checks with atomic Redis SETNX locking to prevent concurrent double-booking of the same seat |
| **Reduced Redis round-trips** | Eliminated redundant Redis GET operations by using lock-based seat holding instead of read-before-write validation |

---

## License

This project is intended for educational and training purposes.
