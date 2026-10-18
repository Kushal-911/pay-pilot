# Session Handoff

## Implemented
- Go/Gin API: `GET /health`, order create/list/detail, PayPal order create, server-side capture.
- PostgreSQL initial migration; Compose service; environment template.
- PayPal OAuth/Orders adapter and focused fake-transport tests.
- OpenAPI, Postman starter collection, README, roadmap/status/decision records.

## What works
- `GOCACHE=/tmp/paypilot-go-cache go test ./...`
- `docker compose config -q`
- The adapter unit test verifies OAuth use and conversion of 2500 USD minor units to `25.00`.

## Run commands
1. Copy `.env.example` to `.env` and add Sandbox credentials.
2. Run `docker compose up --build`.
3. Use `postman/PayPilot.postman_collection.json` or follow `README.md` to create an order, open the approval URL, and capture.

## Broken or incomplete
- No PostgreSQL/Docker runtime flow or live PayPal Sandbox checkout has been verified; local Docker daemon was unreachable.
- PayPal idempotency, unknown-outcome persistence, webhooks, AI, and UI are not implemented.

## Best next task
Run the Compose stack and exercise the Day 1 API against PostgreSQL. Then use configured PayPal Sandbox credentials to verify checkout and capture.
