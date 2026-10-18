# Project Status

## Current date and milestone
2026-10-18 — Day 3 investigation backend and Day 4 webhook/reconciliation backend implemented; runtime integration remains unverified.

## Last completed milestone
No full day milestone has met its end-to-end acceptance criterion. Day 1/2 backend foundations and Day 3 investigation API are implemented, but PostgreSQL runtime and real Sandbox verification are outstanding.

## Implementation audit

| Area | Assessment | Evidence and limits |
|---|---|---|
| Application foundation | Partial | Go/Gin service, Docker Compose, SQL migrations, order and payment routes exist. No frontend screens exist. |
| PayPal Sandbox | Partial | OAuth, order creation, capture, and refresh use PayPal REST endpoints. No Sandbox credential or buyer flow has been run. |
| Order/payment persistence | Partial | PostgreSQL schema has orders and payment attempts, constraints and indexes. Database migrations and endpoints have not been executed against a live DB. |
| Payment lifecycle | Partial | Order statuses pending/paid/failed/review_required; attempt statuses created/capturing/unknown/completed/failed. Capture is guarded and uncertain outcomes require refresh. Runtime transitions unverified. |
| Idempotency | Partial | Mutation middleware stores response/status by request hash and rejects mismatched reuse. PostgreSQL integration test is gated by `TEST_DATABASE_URL` and has not run. |
| Webhooks | Partial | PayPal postback signature verification, event ID deduplication, minimal event storage, and transactional event processing are implemented. Fake verifier and transition tests pass; PayPal and PostgreSQL runtime not verified. |
| Reconciliation | Partial | Provider refresh, order/attempt comparisons, and monotonic capture-event transitions exist. Event sequence coverage is limited; database flow unverified. |
| AI investigation | Partial | Four scoped read-only functions, OpenAI Responses tool loop, ID validation, deterministic fallback, and persisted report API are implemented. No API key found/configured; live model and database flow unverified. |
| Frontend | Missing | No Next.js app or merchant UI exists. |
| Testing/deployment | Partial | Go tests/vet, Compose config, and Postman JSON validation pass. No DB/Sandbox integration test or deployment; local Docker daemon is unavailable. |

## Working features
- Go/Gin API with database health check and order create/list/detail endpoints.
- PostgreSQL schema for orders and payment attempts, initialized by Docker Compose.
- PayPal adapter for OAuth, REST order creation, and capture; backend persists the returned provider identifiers and marks paid only after a completed capture response.
- Mutation idempotency for order creation, checkout creation, and capture; key/request mismatches conflict and completed responses replay.
- Deterministic PayPal request IDs; capture transitions through `capturing` and `unknown`, and refresh reconciles authoritative provider order/capture state before another capture is allowed.
- Order payment-attempt history endpoint.
- Investigation create/get/list API with four order-scoped read-only tools, deterministic discrepancy rules, optional OpenAI function calling, evidence reference validation, and fallback reports.
- PayPal webhook postback signature verification, event ID deduplication, sanitized persisted metadata, and monotonic capture event transitions.
- API and provider tests using fake HTTP transports; a PostgreSQL-gated idempotency integration test is available.
- OpenAPI and Postman starter collection.
- Provider adapter tests using a local fake HTTP transport; PostgreSQL-gated idempotency integration test.

## Known bugs / limitations
- API has not yet been exercised against a running PostgreSQL instance.
- PayPal Sandbox and OpenAI provider calls have not been exercised; credentials are unavailable.
- Idempotency, investigation persistence, and payment transitions have not been run against PostgreSQL.
- Full audit/approval workflow, authorization, and UI remain incomplete.
- API has no merchant authentication/tenant isolation and is local-development only.
- Amount conversion currently assumes two-decimal currencies.

## Current blockers
- PayPal Sandbox credentials and buyer account are not available in the repository/environment, so a real checkout cannot yet be verified.
- Docker CLI is present, but the local Docker daemon is not running/reachable; the Compose stack could not be started for integration checks.

## Files changed
- `PROJECT_STATUS.md`, `ROADMAP.md`, `DECISIONS.md`, `HANDOFF.md`
- `cmd/api/main.go`, `internal/investigations/*.go`, `internal/payments/handlers.go`, `internal/payments/idempotency.go`, `internal/payments/idempotency_test.go`, `internal/payments/webhooks*.go`, `internal/paypal/client.go`, `internal/paypal/client_test.go`
- `migrations/003_investigations.sql`, `docker-compose.yml`, `.env.example`
- `README.md`, `openapi.yaml`, `postman/PayPilot.postman_collection.json`

## Tests run and results
- `go version`: Go 1.26.3 available.
- `node --version`: v22.14.0 available.
- `docker --version`: Docker 29.4.0 available.
- `GOCACHE=/tmp/paypilot-go-cache go test -count=1 ./...`: passed this session, including deterministic investigation, evidence validation, fallback, and fake Responses API tool-loop tests.
- This test run also passed webhook request-body preservation, amount parsing, and out-of-order capture transition tests using fake provider responses. The PostgreSQL integration test is skipped because `TEST_DATABASE_URL` is not set.
- `GOCACHE=/tmp/paypilot-go-cache go vet ./...`: passed this session.
- `docker compose config -q`: passed this session with all three migrations and optional AI environment variables.
- `python3 -m json.tool postman/PayPilot.postman_collection.json`: passed this session.
- `git diff --check`: passed this session.
- `docker compose up --build -d`: blocked because the local Docker daemon is not running/reachable.
- No real database or PayPal Sandbox integration test has been run.

## Required environment variables
Required: `DATABASE_URL`. Optional: `PORT`, `PAYPAL_BASE_URL`, `PAYPAL_CLIENT_ID`, `PAYPAL_CLIENT_SECRET`, `PAYPAL_RETURN_URL`, `PAYPAL_CANCEL_URL`, `PAYPAL_WEBHOOK_ID`, `OPENAI_API_KEY`, `OPENAI_MODEL`, `OPENAI_BASE_URL`. Names only; no values stored here.

## Next three tasks
1. Start Docker/PostgreSQL when available; apply migrations 001–003 and run database-gated idempotency plus order/investigation/webhook flows.
2. Build a minimal merchant UI for the order → attempts → event timeline → investigation evidence journey.
3. Configure Sandbox/OpenAI credentials and verify live provider/model paths; define authorization before adding approval controls.

## Git and deadline
- Branch: `main`
- Latest relevant implementation commit: pending this session's investigation/webhook commit.
- Remaining time: five planned workdays per updated direction; external deadline not supplied.
- Scope cuts: defer AG Grid, APIMatic, Render, analytics, Elastic, and advanced risk scoring until the P0 flow is stable.
