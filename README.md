# Mock Bank API

Mock Bank API is a Go REST service for managing users, authenticated sessions, bank
accounts, and transfers. It uses PostgreSQL for persistence and exposes JSON
endpoints on port `8080` by default.

The service is designed as a small, testable banking backend:

- Gin handles HTTP routing, request binding, validation, and middleware.
- PostgreSQL stores users, sessions, accounts, ledger entries, and transfers.
- sqlc generates type-safe Go code from SQL queries.
- Transfers run as database transactions so the transfer record, ledger entries,
  and account balances are updated together.
- PASETO v2 tokens authenticate API requests. Refresh sessions are persisted in
  PostgreSQL and can be checked or blocked before issuing a new access token.

> This is a development/demo project. The example credentials and token key are
> not suitable for production use.

## Features

- User registration with username, password, full name, and email validation.
- Password hashing and password verification using `golang.org/x/crypto`.
- Login with short-lived access tokens and longer-lived refresh tokens.
- Refresh-token validation against the stored session.
- Authenticated account creation, lookup, pagination, and deletion.
- Supported account currencies: `CAD`, `EUR`, and `USD`.
- One account per user and currency.
- Authenticated transfers between accounts in the same currency.
- Atomic transfer processing with balance updates and debit/credit entries.
- Database schema versioning with `golang-migrate`.
- Mockable storage interface for API unit tests.
- Docker and Docker Compose configurations for local execution.

## API

All request and response bodies use JSON. Protected endpoints require:

```text
Authorization: Bearer <access-token>
```

| Method   | Path                              | Auth | Purpose                                            |
| -------- | --------------------------------- | ---- | -------------------------------------------------- |
| `POST`   | `/users`                          | No   | Create a user                                      |
| `POST`   | `/users/login`                    | No   | Authenticate and create a refresh session          |
| `POST`   | `/tokens/renew_access`            | No   | Exchange a valid refresh token for an access token |
| `POST`   | `/accounts`                       | Yes  | Create an account                                  |
| `GET`    | `/accounts/:id`                   | Yes  | Get an account owned by the authenticated user     |
| `GET`    | `/accounts?page_id=1&page_size=5` | Yes  | List the authenticated user's accounts             |
| `DELETE` | `/accounts/:id`                   | Yes  | Delete an account                                  |
| `POST`   | `/transfers`                      | Yes  | Transfer money between same-currency accounts      |

Account list requests use one-based `page_id`; `page_size` must be between `5`
and `10`. Transfer amounts are positive integer values in the smallest unit
represented by the application (for example, cents if the client uses cents).

## Design

The application follows a small layered design:

1. `main.go` loads configuration, opens the PostgreSQL connection, creates the
   storage implementation, and starts the HTTP server.
2. The `api` package defines routes, request/response types, validation,
   authentication middleware, and handlers.
3. The `db/sqlc` package provides generated query methods and the `Store`
   interface. `db/mock` supplies generated mocks for handler tests.
4. The store's `TransferTx` method coordinates a transaction that creates the
   transfer and entries, then updates both balances. Account IDs are updated in
   a consistent order to reduce deadlock risk.
5. The `token` package abstracts token creation and verification. The current
   implementation uses encrypted PASETO v2 tokens with a 32-byte symmetric key.
6. The `util` package owns configuration loading, password helpers, random
   values, and currency validation.

The database schema is applied incrementally from `db/migration`. Users own
accounts, accounts own ledger entries, and transfers reference both the source
and destination accounts. Foreign keys and unique constraints enforce the core
relationships in PostgreSQL.

## Repository structure

```text
.
├── api/                 HTTP handlers, routes, middleware, and API tests
├── db/
│   ├── migration/       Versioned PostgreSQL schema migrations
│   ├── mock/            Generated storage mocks
│   ├── query/           SQL queries used by sqlc
│   └── sqlc/            Generated models, queries, and store implementation
├── token/               Token interface, payloads, and PASETO implementation
├── util/                Configuration, password, currency, and random helpers
├── .github/workflows/   GitHub Actions CI configuration
├── Dockerfile           Multi-stage container build
├── docker-compose.yaml  PostgreSQL and API services
├── Makefile             Database, generation, server, and test commands
├── app.env              Local runtime configuration
├── .env.example         Configuration template
└── main.go              Application entry point
```

Generated files under `db/sqlc` and `db/mock` should be regenerated from the
source SQL and interfaces rather than edited by hand.

## Prerequisites

For local development, install:

- Go `1.26` or newer compatible with the module's Go version
- PostgreSQL 12, or Docker
- `golang-migrate`
- `sqlc` and `mockgen` if regenerating generated code

## Run locally

1. Start PostgreSQL and create the application database:

   ```sh
   make postgres
   make createdb
   ```

2. Set `DB_SOURCE` in `app.env` to a valid PostgreSQL connection string. For
   the local database created above, the value is:

   ```text
   postgres://root:secret@localhost:5432/bank-api?sslmode=disable
   ```

   Keep `TOKEN_SYMMETRIC_KEY` at exactly 32 characters for local development,
   and replace it with a securely managed secret outside development.

3. Apply migrations and start the server:

   ```sh
   make migrateup
   make server
   ```

   The API is then available at `http://localhost:8080`.

To run the application stack in containers instead:

```sh
docker compose up --build
```

The Compose setup starts PostgreSQL, waits for its TCP port, applies migrations
from the container entrypoint, and starts the API. Stop it with:

```sh
docker compose down
```

## Database and code generation

Create a migration with:

```sh
migrate create -ext sql -dir db/migration -seq <migration_name>
```

Apply or roll back migrations with the corresponding `make` targets:

```sh
make migrateup
make migrateup1
make migratedown1
make migratedown
```

After changing a query or storage interface, regenerate code:

```sh
make sqlc
make mock
```

## Testing

Run the complete Go test suite with coverage output:

```sh
make test
```

This runs `go test -v -cover ./...`, covering:

- API handler and middleware behavior using generated storage mocks.
- Database/sqlc query and transaction behavior.
- Token creation and verification for PASETO and JWT implementations.
- Password hashing and validation utilities.

Tests that exercise the real database expect PostgreSQL to be available and
the schema to have been migrated. The same migration setup used locally can be
used before running tests:

```sh
make postgres
make createdb
make migrateup
make test
```

## Continuous integration

CI is defined in
[`.github/workflows/ci.yaml`](./.github/workflows/ci.yaml) and runs on pushes
to `master` and `ft/docker`, plus pull requests targeting `master`.

The `test` job:

1. Runs on `ubuntu-latest`.
2. Starts a PostgreSQL 12 service with the `bank-api` database and waits for
   its health check to pass.
3. Installs Go 1.26.
4. Checks out the repository.
5. Copies `.env.example` to `app.env`.
6. Installs `golang-migrate` 4.19.1.
7. Applies all migrations with `make migrateup`.
8. Runs `make test`.

A pull request is therefore expected to pass both migration setup and the full
Go test suite before it is considered healthy. CI does not currently publish a
container image or deploy the service.

## Useful commands

| Command                     | Description                           |
| --------------------------- | ------------------------------------- |
| `make server`               | Run the API locally                   |
| `make test`                 | Run all Go tests with coverage        |
| `make postgres`             | Start the local PostgreSQL container  |
| `make createdb`             | Create the `bank-api` database        |
| `make dropdb`               | Drop the `bank-api` database          |
| `make migrateup`            | Apply all migrations                  |
| `make migratedown`          | Roll back all migrations              |
| `make sqlc`                 | Regenerate sqlc code                  |
| `make mock`                 | Regenerate the storage mock           |
| `docker compose up --build` | Build and run PostgreSQL plus the API |
