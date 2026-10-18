# Session Handoff

## Implemented
- Go/Gin API: `GET /health`, order create/list/detail, PayPal order create, server-side capture.
- PostgreSQL initial and reliability migrations; Compose service; environment template.
- Mutation idempotency, payment-attempt history, guarded capture transitions, deterministic PayPal request IDs, and provider refresh/reconciliation.
- PayPal OAuth/Orders adapter and fake-transport tests; PostgreSQL-gated idempotency integration test.
- OpenAPI, Postman starter collection, README, roadmap/status/decision records.

## What works
- `GOCACHE=/tmp/paypilot-go-cache go test ./...`
- `GOCACHE=/tmp/paypilot-go-cache go vet ./...`
- `docker compose config -q`
- Adapter tests verify OAuth use, USD minor-unit conversion, stable provider request IDs, and provider refresh decoding.
- The PostgreSQL integration test is skipped unless `TEST_DATABASE_URL` is provided.

## Run commands
1. Copy `.env.example` to `.env` and add Sandbox credentials.
2. Run `docker compose up --build`.
3. Use `postman/PayPilot.postman_collection.json` or follow `README.md` to create an order, open the approval URL, and capture.

## Broken or incomplete
- No PostgreSQL runtime flow or live PayPal Sandbox checkout has been verified; local Docker daemon was unreachable.
- Webhook verification, AI, and UI are not implemented.

## Best next task
Run the PostgreSQL-gated idempotency test and API flow against the Compose stack, then exercise safe capture recovery using a provider fake or Sandbox credentials.
