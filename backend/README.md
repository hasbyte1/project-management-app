# Project Management Backend API

A robust RESTful API backend for project management built with Go, PostgreSQL, Redis, and WebSockets.

## Features

- **RESTful API** with Chi router
- **PostgreSQL** database with clean architecture
- **Redis** for caching and real-time features
- **JWT authentication** with access and refresh tokens
- **WebSocket** support for real-time updates
- **Clean architecture** (handler -> service -> repository)
- **CORS** enabled for frontend integration
- **Graceful shutdown** support
- **Request logging** and error handling

## Tech Stack

- **Go** 1.21+
- **Chi** - HTTP router
- **PostgreSQL** - Primary database
- **Redis** - Caching and pub/sub
- **JWT** - Authentication
- **sqlx** - SQL extensions
- **bcrypt** - Password hashing

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── models/               # Data models and DTOs
│   ├── repository/           # Database layer
│   ├── service/              # Business logic layer
│   ├── handler/              # HTTP handlers
│   ├── middleware/           # HTTP middleware
│   └── websocket/            # WebSocket handlers
├── pkg/
│   ├── database/             # Database utilities
│   ├── redis/                # Redis client
│   ├── jwt/                  # JWT utilities
│   └── validation/           # Validation helpers
├── migrations/               # Database migrations
└── scripts/                  # Utility scripts
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Redis 6+ (optional)
- Make (optional)

### Installation

1. **Clone the repository**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Set up PostgreSQL database**
   ```bash
   createdb project_management
   ```

5. **Run database migrations**
   ```bash
   # Using the SQL schema provided
   psql -d project_management -f ../path/to/schema.sql
   ```

6. **Run the server**
   ```bash
   go run cmd/server/main.go
   # Or using make
   make run
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/auth/me` - Get current user (protected)

### Tasks

- `POST /api/tasks` - Create task
- `GET /api/tasks/{id}` - Get task by ID
- `PATCH /api/tasks/{id}` - Update task
- `DELETE /api/tasks/{id}` - Delete task
- `PATCH /api/tasks/{id}/status` - Update task status
- `GET /api/projects/{projectId}/tasks` - List project tasks
- `GET /api/tasks/{id}/comments` - Get task comments
- `POST /api/tasks/{id}/comments` - Create comment

### Task Statuses

- `GET /api/projects/{projectId}/statuses` - Get project statuses
- `POST /api/projects/{projectId}/statuses` - Create status

## Configuration

Configuration is managed through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `project_management` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `JWT_ACCESS_SECRET` | JWT access token secret | (required in production) |
| `JWT_REFRESH_SECRET` | JWT refresh token secret | (required in production) |
| `JWT_ACCESS_EXPIRY_MINUTES` | Access token expiry | `15` |
| `JWT_REFRESH_EXPIRY_DAYS` | Refresh token expiry | `7` |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `http://localhost:5173` |

## Development

### Running with auto-reload

Install [Air](https://github.com/cosmtrek/air) for hot reloading:

```bash
go install github.com/cosmtrek/air@latest
make watch
```

### Running tests

```bash
make test
```

### Code formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

## Architecture

### Clean Architecture Layers

1. **Handler Layer** (`internal/handler`)
   - HTTP request/response handling
   - Input validation
   - Authentication checks
   - Response formatting

2. **Service Layer** (`internal/service`)
   - Business logic
   - Data transformation
   - Cross-cutting concerns
   - Transaction management

3. **Repository Layer** (`internal/repository`)
   - Database operations
   - Query building
   - Data mapping
   - Transaction support

### Authentication Flow

1. User registers or logs in
2. Server generates JWT access and refresh tokens
3. Access token included in Authorization header for API requests
4. When access token expires, use refresh token to get new access token
5. Middleware validates tokens on protected routes

### Error Handling

- Consistent error responses across all endpoints
- HTTP status codes follow REST conventions
- Detailed error messages in development
- Generic error messages in production

## Security

- Passwords hashed with bcrypt
- JWT tokens for stateless authentication
- CORS protection
- SQL injection prevention with parameterized queries
- Rate limiting (recommended for production)
- HTTPS required in production

## Database Migrations

Migrations should be created for schema changes:

```bash
make migrate-create NAME=add_new_table
make migrate-up
```

## Docker Support

Build and run with Docker:

```bash
make docker-up
```

## Contributing

1. Follow Go coding standards
2. Write tests for new features
3. Update documentation
4. Use conventional commits

## License

MIT License
