# Chirpy (learn-http-servers)

A social media platform API built with Go using a guide from boot.dev.

## Getting Started

### Prerequisites

- Go 1.24 or later
- PostgreSQL 15 or later
- [goose](https://github.com/pressly/goose) for database migrations
- [sqlc](https://sqlc.dev/) for generating Go code from SQL

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="your-secret-key"
POLKA_KEY="your-polka-api-key"
```

### Database Setup

1. Install PostgreSQL (if not already installed):

```sh
brew install postgresql@15
brew services start postgresql@15
```

2. Create a PostgreSQL database named `chirpy`:

```sh
createdb chirpy
```

3. Install database tools:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

4. Set up the database schema:

```sh
goose -dir sql/schema postgres "postgres://username:password@localhost:5432/chirpy?sslmode=disable" up
```

5. Generate Go database code:

```sh
sqlc generate
```

## API Endpoints

### Authentication

#### Create User

- **POST** `/api/users`
- Creates a new user account
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "secretpassword"
  }
  ```
- **Response** `201`: Returns the created user
  ```json
  {
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "is_chirpy_red": false
  }
  ```

#### Login

- **POST** `/api/login`
- Authenticates a user and returns tokens
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "secretpassword"
  }
  ```
- **Response** `200`:
  ```json
  {
    "id": "uuid",
    "email": "user@example.com",
    "token": "jwt_access_token",
    "refresh_token": "refresh_token",
    "is_chirpy_red": false
  }
  ```

#### Refresh Token

- **POST** `/api/refresh`
- Gets a new access token using a refresh token
- **Header**: `Authorization: Bearer <refresh_token>`
- **Response** `200`:
  ```json
  {
    "token": "new_jwt_access_token"
  }
  ```

#### Revoke Token

- **POST** `/api/revoke`
- Revokes a refresh token
- **Header**: `Authorization: Bearer <refresh_token>`
- **Response** `204`: No content

### User Operations

#### Update User

- **PUT** `/api/users`
- Updates the authenticated user's information
- **Header**: `Authorization: Bearer <access_token>`
- **Body**:
  ```json
  {
    "email": "newemail@example.com",
    "password": "newpassword"
  }
  ```
- **Response** `200`: Returns the updated user

### Chirps

#### Create Chirp

- **POST** `/api/chirps`
- Creates a new chirp
- **Header**: `Authorization: Bearer <access_token>`
- **Body**:
  ```json
  {
    "body": "Hello, world!"
  }
  ```
- **Response** `201`: Returns the created chirp

#### Get All Chirps

- **GET** `/api/chirps`
- Returns all chirps
- **Query Parameters**:
  - `author_id` (optional): Filter chirps by author
- **Response** `200`: Returns array of chirps

#### Get Chirp

- **GET** `/api/chirps/{chirpID}`
- Returns a specific chirp
- **Response** `200`: Returns the chirp

#### Delete Chirp

- **DELETE** `/api/chirps/{chirpID}`
- Deletes a chirp (must be the author)
- **Header**: `Authorization: Bearer <access_token>`
- **Response** `204`: No content

### Webhooks

#### Polka Webhook

- **POST** `/api/polka/webhooks`
- Handles Polka payment webhooks
- **Header**: `Authorization: ApiKey <polka_key>`
- **Body**:
  ```json
  {
    "event": "user.upgraded",
    "data": {
      "user_id": "user_uuid"
    }
  }
  ```
- **Response** `204`: No content

### Administration

#### Reset (Development Only)

- **POST** `/admin/reset`
- Resets the database (only available in development)
- **Response** `200`: Confirmation message

### Health Check

#### Readiness Check

- **GET** `/healthz`
- Checks if the server is ready
- **Response** `200`: "OK"

## Error Responses

All endpoints may return these error responses:

- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Authentication

The API uses JWT tokens for authentication:

1. Access tokens are valid for 1 hour
2. Refresh tokens are valid for 60 days
3. Include tokens in the Authorization header as `Bearer <token>`
