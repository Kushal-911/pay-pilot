# Project Status

## Snapshot — 2026-11-05

PayPilot has a working local prototype with a responsive workbench, Go/Gin API, PostgreSQL persistence, deterministic investigation fallback, PayPal Sandbox adapter, verified webhook handler, and guarded capture recovery. The full local stack and core database paths have now been exercised against fresh PostgreSQL 16 databases. A real PayPal Sandbox checkout and live OpenAI investigation have not been run because account credentials are not available.

**Delivery status:** V0 foundation is implemented and locally verified. V1 hackathon MVP is partially complete: the UI and local API/DB demo work, while real provider/model validation and owner-controlled submission assets remain. The V1 workbench accessibility and narrow-screen checks (#22, #23) are complete. Do not describe this prototype as production-ready; it has no merchant authentication or tenant isolation.

## What works

- Embedded workbench with Synthetic Demo and Live API modes, payment queue/detail, order creation, event timeline, investigation evidence, and provider status/recovery actions.
- Go/Gin API with PostgreSQL-backed orders, payment attempts, idempotency records, investigations, and sanitized provider-event history.
- PayPal adapter for OAuth, order creation, capture, refresh, and signature verification. Fake-provider integration tests exercise these contracts; no live PayPal Sandbox request has been made.
- Capture ambiguity persists as `unknown`. A database-backed test proves a provider refresh is required before capture is eligible again.
- Webhook signature rejection, event deduplication, amount/currency mismatch handling, and completed-state monotonicity are covered against PostgreSQL using a deterministic verification transport.
- Liveness (`/livez`) is separate from database readiness (`/readyz`); `/health` remains a compatibility alias for readiness.
- Demo mode uses clearly labeled synthetic behavior and does not call PayPal or move money.

## Verification completed

The following passed on 2026-11-05:

- Go 1.26.3: `TEST_DATABASE_URL='postgres://paypilot:paypilot_test_local@127.0.0.1:55432/paypilot?sslmode=disable' GOCACHE=/tmp/paypilot-go-cache go test ./...` passed against disposable PostgreSQL 16; the DB-gated tests ran rather than skipped. The suite covers exact idempotency replay/conflict/concurrency/stale claim recovery; order → PayPal order → capture → refresh → investigation persistence/read; capture timeout → `unknown` → refresh approval → safe capture retry; webhook deduplication, signature rejection, mismatch handling, and monotonic state; and sanitized/scoped event history.
- `go vet ./...`, `go build ./cmd/api`, `node --check cmd/api/web/app.js`, `git diff --check`, and `docker compose config -q` passed.
- A fresh isolated Compose project built the API image, applied migrations 001–003 to a new PostgreSQL volume, started the API after the database health check, returned 200 from `/livez`, `/readyz`, and `/health`, and returned 201 from a database-backed order creation request. The temporary project, network, and volume were removed afterward.
- Docker Engine 29.4.0; Postgres image `postgres:16-alpine`; Node v22.14.0.
- Workbench browser review covered the synthetic recovery walkthrough, keyboard payment selection, reset cancellation, and a tablet-width layout. After PR #71, viewport checks at 320, 360, 390, and 640px confirmed that the document does not overflow horizontally, the payment detail fits, and the queue scrolls within its own container. A rendered text contrast scan found no text below 4.5:1 at the audited layout. GitHub issues #22 and #23 are closed.

## Reproducing local startup and PostgreSQL integration checks

For the normal clean-machine path, copy `.env.example` to `.env` if configuring provider/model credentials, then run:

```sh
docker compose up --build -d
curl -fsS http://localhost:8080/livez
curl -fsS http://localhost:8080/readyz
```

If ports 5432 or 8080 are already occupied, the Compose host bindings can be changed without editing the file:

```sh
POSTGRES_HOST_PORT=55433 API_HOST_PORT=18082 docker compose up --build -d
curl -fsS http://localhost:18082/livez
curl -fsS http://localhost:18082/readyz
```

Compose applies numbered SQL migrations 001–003 only when initializing a new PostgreSQL data volume. For disposable integration tests, point `TEST_DATABASE_URL` at a local PostgreSQL 16 database; the test helper creates a unique schema, applies the migrations in that schema, and drops it during cleanup:

```sh
TEST_DATABASE_URL='postgres://paypilot:password@127.0.0.1:5432/paypilot?sslmode=disable' go test ./...
```

The idempotency claim timeout is two minutes; completed idempotency responses have no automatic expiry and must be retained for replay safety. Never reset a volume that contains data you need.

## Remaining work and gates

- **External provider proof:** configure PayPal Sandbox business/buyer credentials and a reachable HTTPS webhook URL; complete buyer approval, capture, refresh, and provider webhook flow. No credentials or customer data are present in this repository.
- **AI proof:** configure an approved OpenAI API key/model and verify the live Responses API path; deterministic fallback is available and covered by tests.
- **UI review:** V1 keyboard/accessibility and phone/tablet layout checks are complete (#22, #23). Revisit accessibility during future feature work and broaden automated coverage beyond the completed manual/browser audit.
- **Submission package:** owner must select/approve the license and provide account-owned Devpost/video actions; accuracy checks, final demo recording, and deadline recheck remain before submission.
- **Production scope:** authentication, tenant isolation, authorization, secret management, audit controls, retention policy, operational monitoring, backups, deployment hardening, and the production gate remain future milestones. Do not host real merchant data.

See GitHub issues #4–#7, #21–#23, #25–#29, and the V2/V3 roadmap issues for detailed acceptance criteria. See `ROADMAP.md`, `PRD.md`, `research/paypal-hackathon-dossier/`, and `AGENTS.md` for scope, sources, and repository practices.

## Current code constraints

- PayPal amount formatting currently assumes a two-decimal currency exponent.
- The API has no merchant authentication or tenant isolation and is suitable only for trusted local development.
- Live API aggregate health metrics are unavailable; synthetic metric values are labeled and hidden in Live API mode.
- PostgreSQL-gated integration tests skip unless `TEST_DATABASE_URL` is explicitly provided.

## Git and delivery

All implementation changes use Conventional Commits and are reviewed through PRs. Per `AGENTS.md`, documentation changes are committed with a focused Conventional Commit directly to `main`, after confirming the branch and fast-forward safety. Never force-push.
