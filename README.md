# Real Estate Microservices System

A containerized Golang-based microservices system for managing users, listings, and routing through a public API gateway. Built with full test coverage, Redis rate limiting, and PostgreSQL storage.

---

## Services Overview

| Service         | Description                    | Port  |
|-----------------|--------------------------------|--------|
| **user-service**   | Manages user data               | `6001` |
| **listing-service**| Manages property listings       | `6000` |
| **public-api**     | Forwards and enriches requests  | `6002` |
| **user-db**        | PostgreSQL for user service     | `5432` |
| **listing-db**     | PostgreSQL for listing service  | `5433` |
| **redis**          | Redis used for rate limiting    | `6379` |

---

## Project Structure

```
real-estate-system/
├── docker-compose.yaml
├── user-service/
│   └── .env
├── listing-service/
│   └── .env
├── public-api/
│   └── .env
```

---

## ⚙️ Environment Setup

### 1. Clone the repo
```bash
git clone https://github.com/dije06/real-estate-system.git
cd real-estate-system
```

### 2. Create `.env` files for each service

#### `user-service/.env`
```env
DB_HOST=user-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

#### `listing-service/.env`
```env
DB_HOST=listing-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

#### `public-api/.env`
```env
USER_SERVICE_URL=http://user-service:6001
LISTING_SERVICE_URL=http://listing-service:6000
REDIS_HOST=redis
REDIS_PORT=6379
```

---

## Docker Compose

### Build and start all services:
```bash
docker compose up --build
```

To stop and clean:
```bash
docker compose down -v
```

---

## 🔐 Public API Endpoints

```
POST /public-api/users        -> forward to user-service
POST /public-api/listings     -> forward to listing-service
GET  /public-api/listings     -> forward + enrich with user data
```

---

## Running Tests

From each service directory:
```bash
go test -v -cover ./...
```

To generate and view HTML coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Features

- ✨ Full unit & integration test coverage
- 🐘 PostgreSQL per service with separate DB containers
- 📦 Redis-backed rate limiting via middleware
- 🧼 Clean `.env`-driven configuration
- 🔁 Health checks and service startup reliability