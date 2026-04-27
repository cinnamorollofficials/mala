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

## API Reference & Endpoint Details

### 1. Data Plane (AI Transactions)
These endpoints are OpenAI-compatible and require a `Virtual Key` in the `Authorization: Bearer <key>` header.

#### `POST /v1/chat/completions`
- **Function**: The main endpoint for chat-based AI interactions. 
- **Process**: Validates key -> Scrubs PII -> Selects Provider -> Proxies request -> Calculates cost -> Logs usage.
- **Compatibility**: Supports standard OpenAI request bodies.

#### `GET /v1/models`
- **Function**: Lists all available and active models that the current Virtual Key is permitted to use.
- **Response**: Returns a JSON list of models with their provider information.

#### `POST /v1/embeddings`
- **Function**: Used for generating vector embeddings for RAG (Retrieval-Augmented Generation) workflows.
- **Status**: Currently implemented as a placeholder.

---

### 2. Control Plane (Management)
These endpoints are used by administrators to manage the system and monitor costs.

#### `POST /admin/keys`
- **Function**: Issues a new Virtual Key for an internal team or application.
- **Input**: `name`, `total_budget` (USD), and `expires_in_days`.
- **Output**: Returns the unique `sk-gh-xxx` key.

#### `GET /admin/keys`
- **Function**: Retrieves a list of all active virtual keys, their allocated budgets, and current spending.

#### `PATCH /admin/keys/:id/topup`
- **Function**: Adds additional budget (USD) to an existing Virtual Key.
- **Use Case**: When a team hits their budget limit and needs more credits.

#### `POST /admin/providers`
- **Function**: Registers a new LLM vendor (OpenAI, Gemini, etc.) in the system.
- **Detail**: The API Key provided here is automatically encrypted using AES-256 before being stored in the database.

#### `POST /admin/models`
- **Function**: Configures a specific AI model (e.g., `gpt-4o`) and links it to a provider.
- **Detail**: Set input and output pricing here for automatic cost calculation.

#### `GET /admin/models`
- **Function**: Lists all configured models with their associated providers and pricing details.

#### `PUT /admin/models/:id`
- **Function**: Updates model pricing or active status.

#### `GET /admin/usage/summary`
- **Function**: Provides a high-level analytics dashboard for the current day.
- **Metric**: Returns the total USD spent across all virtual keys today.

#### `GET /api/health`
- **Function**: Basic system health check. Verifies connectivity to the PostgreSQL database.

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
Run the migration:
```bash
psql mala_db < database/migrations/000001_init_schema.up.sql
```

## Running the Application
```bash
# Run in development mode
make run

# Build production binary
make build

# Run performance test
make perf-test
```

## Docker Deployment

You can also run the entire stack (App + Postgres) using Docker Compose:

```bash
# Start the stack
make docker-up

# View logs
make docker-logs

# Stop the stack
make docker-down
```

> [!TIP]
> On the first run, Docker will automatically execute the database migrations found in `database/migrations/`.

## License
[MIT](LICENSE)
