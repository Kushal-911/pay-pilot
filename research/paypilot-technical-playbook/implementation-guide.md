# Implementation and learning sequence

## Recommended delivery sequence (Oct 19 to Nov 12)

### 1. Verify what already works (first)

- Start Postgres and apply all migrations from a clean database.
- Run Go tests/vet and the DB-gated idempotency test.
- Configure PayPal Sandbox buyer/business accounts and a public HTTPS callback or local tunnel for webhook testing.
- Execute create -> buyer approval -> capture -> refresh; capture request/response IDs and verify the merchant DB state.
- Test a timed-out or failure simulation where supported; do not claim provider response ambiguity unless reproduced or clearly simulated.

### 2. Build one complete merchant workflow

- Pick a frontend that is easy to run and document; do not introduce a large architecture migration.
- Orders/exception list, payment detail, attempts, event timeline, investigation, action explanation.
- Connect all screens to existing API endpoints; no buttons that only animate.
- Add loading, empty, network error and AI unavailable states.
- Clearly label demo data.

### 3. Prove invariants

- Payment marked paid only from a verified completed PayPal result or permitted reconciliation rule.
- Unknown outcome always refreshes before another capture.
- Provider events verified, deduplicated, correlated, amount/currency-checked, and processed monotonically.
- Investigation cannot modify payment state; evidence references validated.
- Model outage does not break basic checkout, refresh or event persistence.

### 4. Make a judge-ready package

- Add LICENSE after confirming you own the submitted code and third-party licenses are compatible.
- README clean-clone setup, environment variable names, migration, Sandbox buyer steps, demo data, API/AI disclosure, limitations.
- Postman collection and OpenAPI paths reflect current behavior.
- One-minute quickstart and a known-good test scenario.
- Record multiple demo takes; remove secrets, user data, unlicensed assets and unsupported claims.
- Submit early and keep the test environment functioning for judging.

## Architecture choices and tradeoffs

| Choice | Why | Tradeoff / guardrail |
|---|---|---|
| Keep Go/Gin/PostgreSQL | Existing backend and state logic already exist; rewrite wastes scarce time | No new frontend currently; choose simplest cohesive UI and maintain API boundary |
| Monolith with packages | One process is sufficient; easy local startup and transactional state changes | Maintain package boundaries; don't mix LLM logic into payment handlers |
| Direct PayPal Orders v2 adapter | Makes PayPal use explicit and easy to demonstrate | Must configure Sandbox and test actual flows; provider mock tests are not a substitute |
| Postgres event/attempt records | Durable evidence timeline, uniqueness and transaction constraints | Need database runtime test and data retention policy before production |
| Deterministic reconciliation + optional LLM | Correctness independent of model; explain complex findings with language interface | AI should add measurable usability value, not merely rephrase a simple status |
| No autonomous payment action | Safer prototype and easier trust story | Recovery is guided; any later consequential automation needs authorization, approval, audit and idempotency |
| Optional AG Grid | Could help sort/filter an exception queue and qualify for sponsor award if actually used well | Do not add if it delays core functional path; no logo-only integration |

## Learning path for a new engineer

1. Read `README.md`, then run `go test ./...`.
2. Trace one handler: `CreatePayPalOrder` or `CapturePayPalOrder` from Gin route to DB and PayPal client.
3. Read `internal/payments/idempotency.go`; understand key reuse and request hash.
4. Read `internal/payments/webhooks.go`; follow raw body, verification, event insert, transaction correlation, transition rules.
5. Read investigation interfaces and tool validation; note why the model cannot read arbitrary records.
6. Read migrations as data contracts; map each foreign key and uniqueness constraint to the business invariant.
7. Use PayPal's official Sandbox docs and test buyer flow. Never use production credentials.
8. Trace an AI tool call from model request -> local tool execution -> function result -> final response; inspect a failure/fallback path.

## Final verification checklist

- [ ] `go test -count=1 ./...` passes; DB-backed tests actually run (not silently skipped).
- [ ] `go vet ./...` passes.
- [ ] Clean `docker compose up --build` and `GET /health` succeeds.
- [ ] Real PayPal Sandbox create/approval/capture/refresh succeeded during this project period.
- [ ] Verified webhook tested through supported real/simulator method; document simulator limitations.
- [ ] UI displays provider and merchant status separately and timeline chronologically.
- [ ] Payment action policy is covered by deterministic tests.
- [ ] AI report facts link to correct persisted records; invalid references fail closed.
- [ ] AI-disabled mode works and is shown in demo.
- [ ] Demo claims match the recorded runtime behavior.
- [ ] Public repo contains no `.env`, credentials, logs with customer data, or private artifacts.
- [ ] License, installation docs, text submission, and public <3-minute YouTube demo complete.
