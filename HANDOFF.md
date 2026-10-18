# Session Handoff

## Implemented
- Go/Gin API: `GET /health`, order create/list/detail, PayPal order create, server-side capture.
- PostgreSQL initial and reliability migrations; Compose service; environment template.
- Mutation idempotency, payment-attempt history, guarded capture transitions, deterministic PayPal request IDs, and provider refresh/reconciliation.
- Evidence-backed investigation API with four scoped read-only tools, optional OpenAI Responses function-calling adapter, deterministic fallback, and evidence ID validation.
- PayPal webhook postback signature verification, deduplicated event persistence, sanitized event metadata, and guarded capture event processing.
- PayPal OAuth/Orders adapter and fake-transport tests; PostgreSQL-gated idempotency integration test.
- OpenAPI, Postman starter collection, README, roadmap/status/decision records.

## What works
- `GOCACHE=/tmp/paypilot-go-cache go test ./...`
- `GOCACHE=/tmp/paypilot-go-cache go vet ./...`
- `docker compose config -q`
- `python3 -m json.tool postman/PayPilot.postman_collection.json`
- Tests verify OAuth, minor-unit conversion, stable request IDs, provider refresh, OpenAI tool calling, evidence validation/fallback, webhook payload preservation, amount parsing, and monotonic event decisions.
- The PostgreSQL integration test is skipped unless `TEST_DATABASE_URL` is provided.

## Run commands
1. Copy `.env.example` to `.env` and add Sandbox credentials.
2. Run `docker compose up --build`.
3. Use `postman/PayPilot.postman_collection.json` or follow `README.md` to create an order, open the approval URL, and capture.

## Broken or incomplete
- No PostgreSQL runtime flow or live PayPal Sandbox checkout has been verified; local Docker daemon was unreachable.
- UI, authorization, approval/rejection, and general audit log are not implemented.
- OpenAI and PayPal runtime credentials are unavailable; those external paths are not live-verified.
- API currently has no merchant authentication/tenant isolation and should remain local-only.

## Best next task
Build the minimal merchant UI for the order → attempts → event timeline → investigation evidence journey. Then run DB-gated persistence flows when Docker is available.
