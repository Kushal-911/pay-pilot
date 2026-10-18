# Project Status

## Current milestone
Day 2 — Payment lifecycle and reliability. Implementation is in place; PostgreSQL-backed replay/transition behavior still needs execution against a running database.

## Last completed milestone
Day 1 backend foundation and PayPal checkout adapter (implementation complete; Sandbox checkout not externally verified).

## Working features
- Go/Gin API with database health check and order create/list/detail endpoints.
- PostgreSQL schema for orders and payment attempts, initialized by Docker Compose.
- PayPal adapter for OAuth, REST order creation, and capture; backend persists the returned provider identifiers and marks paid only after a completed capture response.
- Mutation idempotency for order creation, checkout creation, and capture; key/request mismatches conflict and completed responses replay.
- Deterministic PayPal request IDs; capture transitions through `capturing` and `unknown`, and refresh reconciles authoritative provider order/capture state before another capture is allowed.
- Order payment-attempt history endpoint.
- OpenAPI and Postman starter collection.
- Provider adapter tests using a local fake HTTP transport; PostgreSQL-gated idempotency integration test.

## Known bugs / limitations
- API has not yet been exercised against a running PostgreSQL instance.
- PayPal Sandbox flow has not been exercised; credentials and buyer account are unavailable.
- Idempotency and provider refresh code is implemented but has not been executed against PostgreSQL.
- Webhook verification, request IDs, and reconciliation remain for Day 4.
- Amount conversion currently assumes two-decimal currencies.

## Current blockers
- PayPal Sandbox credentials and buyer account are not available in the repository/environment, so a real checkout cannot yet be verified.
- Docker CLI is present, but the local Docker daemon is not running/reachable; the Compose stack could not be started for integration checks.

## Files changed
- `PROJECT_STATUS.md`, `ROADMAP.md`, `DECISIONS.md`, `HANDOFF.md`
- `cmd/api/main.go`, `internal/payments/handlers.go`, `internal/payments/idempotency.go`, `internal/payments/idempotency_test.go`, `internal/paypal/client.go`, `internal/paypal/client_test.go`
- `migrations/002_reliability.sql`, `docker-compose.yml`
- `README.md`, `openapi.yaml`, `postman/PayPilot.postman_collection.json`

## Tests run and results
- `go version`: Go 1.26.3 available.
- `node --version`: v22.14.0 available.
- `docker --version`: Docker 29.4.0 available.
- `GOCACHE=/tmp/paypilot-go-cache go test ./...`: passed; PayPal adapter fake-transport tests passed, other packages compile.
- `GOCACHE=/tmp/paypilot-go-cache go test ./...`: passed on Day 2; PostgreSQL integration test skipped because `TEST_DATABASE_URL` is not set.
- `GOCACHE=/tmp/paypilot-go-cache go vet ./...`: passed.
- `docker compose config -q`: passed.
- `python3 -m json.tool postman/PayPilot.postman_collection.json`: passed.
- `git diff --check`: passed.
- `docker compose up --build -d`: blocked because the local Docker daemon is not running/reachable.
- No real database or PayPal Sandbox integration test has been run.

## Required environment variables
Planned: `DATABASE_URL`, `PORT`, `PAYPAL_BASE_URL`, `PAYPAL_CLIENT_ID`, `PAYPAL_CLIENT_SECRET`, `PAYPAL_RETURN_URL`, `PAYPAL_CANCEL_URL`.

## Next three tasks
1. Start Docker/PostgreSQL, apply migration 002, and run `TEST_DATABASE_URL=... go test ./...` to verify idempotency replay and mismatch behavior.
2. Exercise capture failure, refresh-to-approved, safe retry, and refresh-to-completed with a provider fake plus PostgreSQL.
3. Configure PayPal Sandbox credentials/accounts and verify checkout approval, capture, and refresh end to end.

## Git and deadline
- Branch: `main`
- Latest relevant implementation commit: `43dc1a5` (`feat: add payment idempotency and recovery states`).
- Remaining time: approximately 7 calendar days from 2026-10-18; hackathon deadline not supplied.
- Scope cuts: defer AG Grid, APIMatic, Render, analytics, Elastic, and advanced risk scoring until the P0 flow is stable.
