# Project Status

## Current milestone
Day 1 — Foundation and PayPal checkout. Initial backend slice is implemented; next is exercising it against PostgreSQL and configured Sandbox credentials.

## Last completed milestone
Repository inspection (2026-10-18). The repository was empty on branch `main` with no commits.

## Working features
- Go/Gin API with database health check and order create/list/detail endpoints.
- PostgreSQL schema for orders and payment attempts, initialized by Docker Compose.
- PayPal adapter for OAuth, REST order creation, and capture; backend persists the returned provider identifiers and marks paid only after a completed capture response.
- OpenAPI and Postman starter collection.
- Provider adapter unit tests using a local fake HTTP transport.

## Known bugs / limitations
- API has not yet been exercised against a running PostgreSQL instance.
- PayPal Sandbox flow has not been exercised; credentials and buyer account are unavailable.
- Request idempotency, provider refresh, webhook verification, request IDs, and unknown capture outcome persistence remain for Day 2/4.
- Amount conversion currently assumes two-decimal currencies.

## Current blockers
- PayPal Sandbox credentials and buyer account are not available in the repository/environment, so a real checkout cannot yet be verified.
- Docker CLI is present, but the local Docker daemon is not running/reachable; the Compose stack could not be started for integration checks.

## Files changed
- `PROJECT_STATUS.md`, `ROADMAP.md`, `DECISIONS.md`, `HANDOFF.md`
- `go.mod`, `go.sum`, `cmd/api/main.go`, `internal/payments/handlers.go`, `internal/paypal/client.go`, `internal/paypal/client_test.go`
- `migrations/001_initial.sql`, `docker-compose.yml`, `Dockerfile`, `.env.example`, `.gitignore`
- `README.md`, `openapi.yaml`, `postman/PayPilot.postman_collection.json`

## Tests run and results
- `go version`: Go 1.26.3 available.
- `node --version`: v22.14.0 available.
- `docker --version`: Docker 29.4.0 available.
- `GOCACHE=/tmp/paypilot-go-cache go test ./...`: passed; PayPal adapter fake-transport tests passed, other packages compile.
- `docker compose config -q`: passed.
- `docker compose up --build -d`: blocked because the local Docker daemon is not running/reachable.
- No real database or PayPal Sandbox integration test has been run.

## Required environment variables
Planned: `DATABASE_URL`, `PORT`, `PAYPAL_BASE_URL`, `PAYPAL_CLIENT_ID`, `PAYPAL_CLIENT_SECRET`, `PAYPAL_RETURN_URL`, `PAYPAL_CANCEL_URL`.

## Next three tasks
1. Start the local Docker daemon, run Compose, and exercise health, create/list/detail, and database constraints.
2. Configure PayPal Sandbox credentials/accounts and verify checkout approval plus server-side capture end to end.
3. Add idempotency and payment state-transition tests; persist unknown provider outcomes before enabling retries.

## Git and deadline
- Branch: `main`
- Latest relevant commit: none (repository has no commits).
- Remaining time: approximately 7 calendar days from 2026-10-18; hackathon deadline not supplied.
- Scope cuts: defer AG Grid, APIMatic, Render, analytics, Elastic, and advanced risk scoring until the P0 flow is stable.
