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

## 2026-10-18 — Deterministic evidence first, optional model explanation second
- Always load order, attempt, and matching event records and run deterministic comparisons. The optional OpenAI Responses adapter receives only four scoped, read-only functions; recommendations remain a fixed enum and deterministic payment-state checks override the model.
- Model output is accepted only when its report fields validate and every cited evidence ID resolves to a record in the loaded set. Missing credentials, provider failures, and invalid model output preserve and return a deterministic report.
- Rationale: investigation remains useful without an LLM account and the model cannot claim a payment outcome or perform a payment mutation.

## 2026-10-18 — Verify PayPal webhooks through PayPal postback verification
- Send the original raw event bytes and PayPal transmission headers to PayPal's `verify-webhook-signature` endpoint using the configured webhook ID. Persist only selected event metadata, not the full event payload.
- Deduplicate by provider event ID and apply capture transitions in one database transaction; completed capture state is never downgraded by later pending/denied events.
- Rationale: use the provider-supported verification path and keep event processing auditable while limiting retained payload data.
