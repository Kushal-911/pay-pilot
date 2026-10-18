# Architecture Decisions

## 2026-10-18 — Modular monolith and required stack
- Use Go with Gin, pgx, and PostgreSQL; organize provider, payment, HTTP, and persistence concerns as modules in one service.
- Rationale: matches the requested stack while keeping a one-developer, seven-day delivery operable.

## 2026-10-18 — Provider adapter owns remote calls
- PayPal REST operations live behind a Go adapter. Backend payment logic validates state and persists provider outcomes; browser redirects and AI text are never authoritative.
- Rationale: isolate provider details and keep payment mutations under deterministic backend control.

## 2026-10-18 — Build incrementally, defer optional integrations
- Start with persisted orders and checkout; defer sponsor/analytics features until the P0 path works.
- Rationale: prioritize a reliable, reviewable payment flow within the deadline.

## 2026-10-18 — Persist idempotency responses and deterministic provider request IDs
- Require an `Idempotency-Key` on mutation endpoints, hash method/path/body, replay completed JSON responses, and conflict on key reuse with a different request.
- Derive PayPal request IDs from internal order/payment IDs so retries map to the same provider operation.
- Rationale: protect both the API boundary and PayPal side effects against duplicate delivery.

## 2026-10-18 — Unknown capture outcomes require provider refresh
- Persist `capturing` before contacting PayPal and `unknown` when the response cannot confirm an outcome. Only a provider refresh can reconcile `unknown` back to `created` or forward to `completed`.
- Rationale: a timeout cannot be treated as a failed capture; retry eligibility comes from current authoritative provider state.
