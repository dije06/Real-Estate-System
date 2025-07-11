# Real Estate Microservices System

A containerized microservices-based backend system written in Go (Golang) using the Echo framework. This project manages users, property listings, and a public-facing API layer through isolated services. It fulfills the full requirements of the Backend Tech Challenge.

## Project Structure

```
real-estate-system/
├── user-service/          - Handles user CRUD operations  
├── listing-service/       - Handles property listing creation and retrieval  
├── public-api/            - Public-facing API layer (gateway)  
├── docker-compose.yaml    - Orchestrates all services and databases  
├── .env.example           - Environment variable template  
└── README.md              - This file
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- make (optional, for simplified CLI usage)

### Quickstart with Docker Compose

```bash
# Clone the repository
git clone https://github.com/dije06/Real-Estate-System.git
cd Real-Estate-System

# Copy env templates
cp .env.example .env

# Start all services
docker-compose up --build
```

## Microservices Overview

### 1. User Service (`localhost:6001`)

Manages users.

- `GET /users`: Paginated list of users  
- `GET /users/:id`: Retrieve a user by ID  
- `POST /users`: Create a user using `application/x-www-form-urlencoded`

### 2. Listing Service (`localhost:6000`)

Manages listings.

- `GET /listings`: Paginated listings, optional filter by `user_id`  
- `POST /listings`: Create a listing using `application/x-www-form-urlencoded`

### 3. Public API (`localhost:6002`)

Gateway for frontend/mobile clients.

- `GET /public-api/listings`: Listings with user detail  
- `POST /public-api/users`: Create user (JSON)  
- `POST /public-api/listings`: Create listing (JSON)

## Example API Calls

### Create User (Internal Service)

```bash
curl -X POST http://localhost:6001/users \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=John Doe"
```

### Create Listing (Internal Service)

```bash
curl -X POST http://localhost:6000/listings \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_id=1" \
  -d "listing_type=rent" \
  -d "price=5000"
```

### Public API - Create User

```bash
curl -X POST http://localhost:6002/public-api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe"}'
```

### Public API - Create Listing

```bash
curl -X POST http://localhost:6002/public-api/listings \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "listing_type": "rent", "price": 6000}'
```

### Architecture

| Component        | Description                                       |
|------------------|---------------------------------------------------|
| `user-service`   | Handles users and uses its own PostgreSQL DB      |
| `listing-service`| Handles listings and uses its own PostgreSQL DB   |
| `public-api`     | Acts as the only entry point for external clients |
| `redis`          | Rate limiting support for public API              |
| `docker-compose` | Spins up all services with one command            |

### REST API Contract Compliance

| Endpoint                   | Required Content-Type              | Implemented As                        |
|----------------------------|------------------------------------|---------------------------------------|
| `POST /users`              | `application/x-www-form-urlencoded`| Uses `c.FormValue("name")`            |
| `POST /listings`           | `application/x-www-form-urlencoded`|  Uses `c.FormValue(...)`              |
| `POST /public-api/users`   | `application/json`                 |  Uses `c.Bind()`                      |
| `POST /public-api/listings`| `application/json`                 |  Uses `c.Bind()`                      |

## Testing

- Use the included Postman collection (`RealEstateAPI.postman_collection.json`) to test endpoints  
- Unit and integration tests are located under each service's `tests/` directory

## Acknowledgements

Built with love to demonstrate clean microservices design using Go and Echo.
