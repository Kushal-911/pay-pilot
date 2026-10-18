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
