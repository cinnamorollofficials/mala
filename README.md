# Mala LLM Gateway

Mala is a high-performance LLM Gateway built with Go, using the Fiber web framework and PostgreSQL. It serves as a unified entry point for interacting with multiple Large Language Models while providing centralized management, logging, and security.

## Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Fiber v2](https://gofiber.io/)
- **Database Driver**: [pgx v5](https://github.com/jackc/pgx)
- **Query Builder**: [Squirrel](https://github.com/Masterminds/squirrel)
- **Environment Management**: [godotenv](https://github.com/joho/godotenv)

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Make (optional, but recommended)

## Getting Started

### 1. Clone the repository

```bash
git clone git@github.com:cinnamorollofficials/mala.git
cd mala
```

### 2. Configure environment variables

Create a `.env` file in the root directory (you can use the existing one as a template):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=mala_db
DB_SSLMODE=disable

PORT=3000
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Database Setup

Ensure you have created the database specified in your `.env` file:

```sql
CREATE DATABASE mala_db;
```

## Running the Application

Using **Make**:

```bash
# Run the application
make run

# Build the binary
make build

# Run tests
make test
```

Without **Make**:

```bash
go run main.go
```

## Project Structure

```text
mala/
├── cmd/                # Entry points for the application
├── internal/           # Private application and library code
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Database models and DTOs
│   ├── routes/         # Route definitions
│   └── service/        # Business logic
├── pkg/                # Public library code (can be used by other projects)
│   └── database/       # Database connection and utilities
├── .env                # Environment variables (git-ignored)
├── .gitignore          # Git ignore rules
├── Makefile            # Build and development commands
└── main.go             # Application entry point
```

## API Endpoints

| Method | Endpoint      | Description                           |
| :----- | :------------ | :------------------------------------ |
| `GET`  | `/api/health` | Check application and database status |

## Author

Initial project by [cinnamorollofficials](https://github.com/cinnamorollofficials)

## License

[MIT](LICENSE)
