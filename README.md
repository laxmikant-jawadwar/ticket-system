# Ticket System API

A simple REST API built with Golang for managing user-created support tickets.

Users can register, login, create tickets, view their own tickets, and update the status of their own tickets.

The project focuses on JWT authentication, ownership-based authorization, PostgreSQL persistence, Dockerization, and deployment.

---

## Features

- User registration
- User login
- JWT-based authentication
- Password hashing using bcrypt
- Create tickets
- View only the authenticated user's tickets
- View a specific own ticket
- Update status of own tickets
- Status transition validation
- Closed tickets cannot be reopened
- PostgreSQL persistence
- Docker support
- Cloud deployment
- Public health check endpoint

---

## Tech Stack

- **Language:** Golang
- **HTTP Server:** Go `net/http`
- **Authentication:** JWT
- **Password Hashing:** bcrypt
- **Database:** PostgreSQL
- **Database Driver:** pgx
- **Containerization:** Docker
- **Deployment:** Render
- **Cloud Database:** Supabase PostgreSQL

---

## Project Structure

```text
ticket-system/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── database.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── health.go
│   │   └── ticket_handler.go
│   ├── middleware/
│   │   └── auth_middleware.go
│   ├── model/
│   │   ├── auth.go
│   │   ├── ticket.go
│   │   └── user.go
│   ├── repository/
│   │   ├── ticket_repository.go
│   │   └── user_repository.go
│   └── service/
│       ├── auth_service.go
│       ├── jwt_service.go
│       └── ticket_service.go
├── .env.example
├── .gitignore
├── Dockerfile
├── README.md
├── go.mod
└── go.sum
```

---

# API Documentation

## Base URL

### Local

```text
http://localhost:8080
```

### Deployed

```text
https://ticket-system-ivf1.onrender.com
```

---

## Authentication

The ticket APIs are protected using JWT authentication.

After successful login, the API returns a JWT token.

The token must be sent with protected requests using:

```http
Authorization: Bearer <token>
```

Passwords are hashed using bcrypt and are never stored as plain text.

---

# API Endpoints

## 1. Health Check

### Request

```http
GET /health
```

### Example

```bash
curl http://localhost:8080/health
```

### Response

```json
{
  "status": "ok"
}
```

---

## 2. Register User

### Request

```http
POST /auth/register
Content-Type: application/json
```

### Request Body

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

### Example

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"John Doe\",\"email\":\"john@example.com\",\"password\":\"password123\"}"
```

A successful registration creates a new user.

---

## 3. Login

### Request

```http
POST /auth/login
Content-Type: application/json
```

### Request Body

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

### Example

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"john@example.com\",\"password\":\"password123\"}"
```

### Response

```json
{
  "token": "<jwt-token>"
}
```

The returned token is required for protected ticket endpoints.

---

## 4. Create Ticket

### Request

```http
POST /tickets
Authorization: Bearer <token>
Content-Type: application/json
```

### Request Body

```json
{
  "title": "Unable to login",
  "description": "I am unable to login to my account."
}
```

### Example

```bash
curl -X POST http://localhost:8080/tickets \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"Unable to login\",\"description\":\"I am unable to login to my account.\"}"
```

New tickets are created with the status:

```text
open
```

---

## 5. List Own Tickets

### Request

```http
GET /tickets
Authorization: Bearer <token>
```

### Example

```bash
curl http://localhost:8080/tickets \
  -H "Authorization: Bearer <token>"
```

Only tickets belonging to the authenticated user are returned.

---

## 6. Get Own Ticket

### Request

```http
GET /tickets/{id}
Authorization: Bearer <token>
```

### Example

```bash
curl http://localhost:8080/tickets/1 \
  -H "Authorization: Bearer <token>"
```

A user can only access a ticket created by that same user.

---

## 7. Update Ticket Status

### Request

```http
PATCH /tickets/{id}/status
Authorization: Bearer <token>
Content-Type: application/json
```

### Request Body

```json
{
  "status": "in_progress"
}
```

### Example

```bash
curl -X PATCH http://localhost:8080/tickets/1/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d "{\"status\":\"in_progress\"}"
```

---

# Ticket Status Flow

Tickets follow a fixed status flow:

```text
open
  |
  v
in_progress
  |
  v
closed
```

### Rules

- New tickets start with `open`.
- `open` can be changed to `in_progress`.
- `in_progress` can be changed to `closed`.
- `closed` tickets cannot be reopened.
- Invalid status values are rejected.
- Users can update only their own tickets.

Valid statuses are:

```text
open
in_progress
closed
```

---

# Setup Instructions

## Prerequisites

The following are required for local development:

- Go
- PostgreSQL

Docker is optional if running the application using Docker.

---

## 1. Clone the Repository

```bash
git clone <your-github-repository-url>
cd ticket-system
```

---

## 2. Configure Environment Variables

Create a `.env` file in the project root.

Use `.env.example` as a reference:

```env
PORT=8080
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/ticket_system
JWT_SECRET=your-secret-key
```

### Environment Variables

| Variable | Description |
|---|---|
| `PORT` | Port on which the API server runs |
| `DATABASE_URL` | PostgreSQL database connection string |
| `JWT_SECRET` | Secret used to sign JWT tokens |


---

# Database Setup

The application uses PostgreSQL.

Create a database named:

```text
ticket_system
```

Then create the required tables:

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

---

# Running Locally

From the project root:

```bash
go run ./cmd/server
```

The API will be available at:

```text
http://localhost:8080
```

Test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

The root endpoint is also available:

```text
http://localhost:8080/
```

and returns:

```text
Welcome to Ticket System API
```

---

# Running with Docker

A Dockerfile is included in the project.

## Build the Docker Image

```bash
docker build -t ticket-system .
```

## Run the Container

```bash
docker run -p 8080:8080 ticket-system
```

For a database-backed container, provide the required environment variables:

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="<database-url>" \
  -e JWT_SECRET="<jwt-secret>" \
  ticket-system
```

Then test:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

# Deployment

The application is deployed using Render.

## Deployed Application

```text
https://ticket-system-ivf1.onrender.com
```

## Public Health Check

```text
https://ticket-system-ivf1.onrender.com/health
```

Expected response:

```json
{
  "status": "ok"
}
```


---

# Cloud Database

The deployed application uses PostgreSQL hosted on Supabase.

The database connection is provided through the `DATABASE_URL` environment variable.

Database credentials and JWT secrets are stored as environment variables and are not committed to the repository.

---

# API Testing Flow

A typical testing flow is:

```text
1. Register a user
       |
       v
2. Login and receive JWT
       |
       v
3. Create a ticket using JWT
       |
       v
4. List the user's tickets
       |
       v
5. Get a ticket by ID
       |
       v
6. Update status
   open -> in_progress -> closed
       |
       v
7. Verify closed ticket cannot be reopened
```

Ownership can also be verified by using a different authenticated user and attempting to access another user's ticket.

---

# HTTP Status Codes

The API uses meaningful HTTP status codes for different outcomes, including:

| Status Code | Meaning |
|---|---|
| `200` | Successful request |
| `201` | Resource created |
| `400` | Invalid request/input |
| `401` | Authentication required/invalid credentials |
| `403` | Unauthorized operation |
| `404` | Resource not found |
| `409` | Resource conflict, such as duplicate email |
| `500` | Internal server error |

---

# Assumptions

- No admin role is implemented because it is not required.
- No ticket assignment flow is implemented.
- No comments module is implemented.
- Each ticket belongs to the user who created it.
- Users can view and update only their own tickets.
- Ticket status transitions are restricted to:
  `open -> in_progress -> closed`.
- Closed tickets cannot be reopened.
- PostgreSQL is used as the persistent data store.
- JWT is used for authentication of protected endpoints.
- Passwords are stored as bcrypt hashes.
- The implementation intentionally keeps the system simple and avoids unnecessary features.

---

# Security

- Passwords are hashed using bcrypt.
- JWT is required for protected ticket endpoints.
- JWT secrets are provided through environment variables.
- Database credentials are provided through environment variables.
- `.env` is excluded from Git using `.gitignore`.
- Ticket ownership is checked before allowing users to view or update tickets.

---

# Note

This project was implemented as a simple backend ticket system with a focus on:

- Correct REST API behavior
- JWT authentication
- Ownership-based authorization
- Secure password storage
- PostgreSQL persistence
- Valid ticket status transitions
- Dockerization
- Working cloud deployment

---

## Thank You

I enjoyed building this Ticket System API and applying backend concepts such as authentication, authorization, REST API design, database persistence, and Docker deployment.

**Built with Go. Focused on simplicity, correctness, and clean backend design.**