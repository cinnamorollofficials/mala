# Mala LLM Gateway

Mala is a high-performance, enterprise-grade LLM Gateway built with Go and Fiber. It serves as a secure, unified entry point for multiple LLM providers (OpenAI, Anthropic, Gemini) with built-in budget management, security scrubbing, and detailed analytics.

## Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Fiber v2](https://gofiber.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [pgx v5](https://github.com/jackc/pgx)
- **Query Builder**: [Squirrel](https://github.com/Masterminds/squirrel)
- **Security**: AES-256 for API key encryption

## Core Features

- **Data Plane (OpenAI-Compatible)**: Drop-in replacement for OpenAI endpoints.
- **Security Chain**:
    - **IP Whitelisting**: Restrict access to trusted internal servers.
    - **Virtual Key Auth**: Manage internal access with virtual keys.
    - **Budget Guard**: Real-time spending enforcement and auto-blocking.
    - **Rate Limiting**: Per-key request throttling.
    - **PII Scrubber**: Automated redaction of sensitive data (Email, NIK, Phone) before sending to providers.
- **Control Plane (Admin)**:
    - **Key Management**: Create and manage virtual keys with specific budgets.
    - **Provider Health**: Monitor uptime of upstream LLM providers.
    - **Analytics**: Cost tracking and usage history.

## How It Works

Mala operates as a transparent proxy between your internal applications and external LLM providers. Here is the lifecycle of a request:

### 1. Request Flow (The Security Chain)
Every request to the Data Plane (`/v1/*`) passes through a rigorous middleware chain:
1.  **IP Filtering**: Ensures the request comes from an authorized internal server IP.
2.  **Virtual Key Validation**: Validates the `Authorization` bearer token against the database.
3.  **Budget Guard**: Checks if the key has enough remaining budget. If the balance is ≤ 0, it returns `402 Payment Required`.
4.  **Rate Limiter**: Enforces Request-Per-Minute (RPM) limits to prevent accidental resource exhaustion.
5.  **PII Scrubber**: Scans the request body (prompt) for sensitive information (Emails, Phone Numbers, NIK) and redacts them using regex patterns.
6.  **Provider Setup**: Identifies the correct upstream provider for the requested model and decrypts the real API key.

### 2. Execution & Logging
- The handler forwards the (scrubbed) request to the upstream provider (e.g., OpenAI).
- Once the response is received, Mala parses the token usage.
- **Asynchronous Logging**: Mala calculates the cost based on the model's price and records the transaction in `usage_logs` without blocking the client response.
- The virtual key's `spent_amount` is updated in real-time.

## Getting Started

### 1. Prerequisites

- Go 1.21+
- PostgreSQL
- Make

### 2. Installation

```bash
git clone git@github.com:cinnamorollofficials/mala.git
cd mala
go mod tidy
```

### 3. Configuration

Create a `.env` file from the example:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=mala_db
DB_SSLMODE=disable

PORT=3000
ENCRYPTION_KEY=your-32-char-encryption-key
ALLOWED_IPS=127.0.0.1,::1
```

### 4. Database Initialization

Run the migration found in `database/migrations/000001_init_schema.up.sql` against your PostgreSQL instance.

## Running the Application

```bash
# Run in development mode
make run

# Build production binary
make build
```

## API Reference

### Data Plane (OpenAI-Compatible)
Requires `Authorization: Bearer <VIRTUAL_KEY>`

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/v1/chat/completions` | Proxy chat requests to configured providers |
| `GET` | `/v1/models` | List active models available to the key |

### Control Plane (Admin)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/admin/keys` | Create a new Virtual Key |
| `GET` | `/admin/keys` | List all Virtual Keys and budgets |
| `PATCH` | `/admin/keys/:id/topup` | Add budget to a key |
| `POST` | `/admin/providers` | Configure a new LLM Provider |
| `GET` | `/admin/usage/summary` | Get total cost analytics for today |

## Project Structure

```text
mala/
├── database/migrations # SQL migration files
├── internal/
│   ├── handlers/       # Data & Control plane handlers
│   ├── middleware/     # Security chain (Auth, PII, Budget, etc.)
│   ├── models/         # Go entities (GORM/SQL compatible)
│   └── routes/         # Router orchestration
├── pkg/
│   ├── database/       # DB connection pool
│   └── utils/          # Encryption & helpers
└── main.go             # Application entry point
```

## License

[MIT](LICENSE)
