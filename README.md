# Real Estate Microservices System

A containerized Golang-based microservices system using Echo web framework for managing users, property listings, and routing through a public API gateway. The system includes PostgreSQL databases for each core service, Redis-based rate limiting, and complete unit + integration test coverage.

---

## Services Overview

Service: `user-service`
- Purpose: Manages user data and exposes user CRUD endpoints
- Port: `6001`
- Database: PostgreSQL (`user-db` on port `5432`)

Service: `listing-service`
- Purpose: Manages real estate property listings
- Port: `6000`
- Database: PostgreSQL (`listing-db` on port `5433`)

Service: `public-api`
- Purpose: Acts as a gateway between frontend clients and internal services, handling request forwarding and enrichment
- Port: `6002`
- Depends on: `user-service`, `listing-service`, `redis`

Service: `redis`
- Purpose: Rate-limiting for public API routes
- Port: `6379`

---

## Project Structure

This project is divided into three main microservices with separate handlers, models, repositories, seeders, tests, and `.env` configurations.

Root structure:

real-estate-system/
- docker-compose.yaml
- README.md
- .gitignore
- postman_collection.json

user-service/
- .env (ignored)
- .env.example
- Dockerfile
- main.go
- go.mod, go.sum
- handlers/
  - user_handler.go
  - tests/user_handler_test.go
- models/
  - user.go
- repository/
  - user_repository.go
  - interfaces/user_repository_interface.go
  - mocks/user_repository_mock.go
  - tests/user_repository_test.go
- seeders/
  - user_seed.go

listing-service/
- .env (ignored)
- .env.example
- Dockerfile
- main.go
- go.mod, go.sum
- handlers/
  - listing_handler.go
  - tests/listing_handler_test.go
- models/
  - listing.go
- repository/
  - listing_repository.go
  - interfaces/listing_repository_interface.go
  - mocks/listing_repository_mock.go
  - tests/listing_repository_test.go
- seeders/
  - listing_seed.go

public-api/
- .env (ignored)
- .env.example
- Dockerfile
- main.go
- go.mod, go.sum
- handlers/
  - public_handler.go
  - tests/public_handler_test.go
- middleware/
  - rate_limiter.go

---

## Environment Setup

Clone the repo:
```bash
git clone https://github.com/dije06/Real-Estate-System.git
cd real-estate-system
```

Set up `.env` files for each service:

user-service/.env:
```
DB_HOST=user-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

listing-service/.env:
```
DB_HOST=listing-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
DB_TIMEZONE=UTC
```

public-api/.env:
```
USER_SERVICE_URL=http://user-service:6001
LISTING_SERVICE_URL=http://listing-service:6000
REDIS_HOST=redis
REDIS_PORT=6379
```

---

## Running the Application

Build and run all services:
```bash
docker compose up --build
```

Stop and remove volumes:
```bash
docker compose down -v
```

---

## Public API Endpoints

- `POST /public-api/users`: Forwards to user-service to create a user
- `POST /public-api/listings`: Forwards to listing-service to create a listing
- `GET /public-api/listings`: Retrieves listings and enriches them with user info

---

## Internal Service API Endpoints

### user-service
- `GET /users?page_num=&page_size=`: Paginated list of users
- `GET /users/:id`: Get single user by ID
- `POST /users`: Create new user

### listing-service
- `GET /listings?page_num=&page_size=&user_id=`: Paginated list of listings (with optional user filter)
- `GET /listings/:id`: Get single listing by ID
- `POST /listings`: Create new listing

---

## Testing

To run tests:
```bash
cd [service-folder]
go test -v -cover ./...
```

Generate and open coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Postman Collection

This repository includes a ready-to-use Postman collection:
- File: `Real Estate System.postman_collection.json`
- Covers all public API routes
- Includes sample requests and expected responses

You can import this file into Postman to test:
- User creation
- Listing creation
- Listing retrieval with user enrichment

---

## Features

- Full coverage for handlers and repositories using `testify` and `sqlmock`
- PostgreSQL databases containerized and isolated per service
- Redis-based rate limiting applied to public API
- Modular folder structure with dedicated tests and seeders
- Environment-variable driven config using Docker `env_file`
- Health checks using `pg_isready` before service startup
- Postman collection for API testing and demonstration

