# PayPilot Roadmap

## Day 1 — Foundation and PayPal checkout (P0)
- [x] Initialize Go/Gin backend and PostgreSQL persistence.
- [x] Add Docker Compose, environment example, health endpoint, and initial SQL migration.
- [x] Implement internal order creation and retrieval.
- [x] Implement PayPal order creation, approval URL response, and server-side capture endpoint.
- [x] Start Postman collection.
- [ ] Exercise the complete Sandbox checkout with configured credentials and buyer approval.
- Acceptance: a Sandbox payment completes and its verified result is persisted and visible. External credentials required.

## Day 2 — Payment lifecycle and reliability (P0)
- [x] Add payment attempts, idempotency, state-transition rules, and payment history/detail.
- [x] Model unknown provider outcomes and safe refresh before consequential retries.
- [x] Add focused adapter tests and a PostgreSQL-gated idempotency integration test.
- [ ] Run PostgreSQL integration and state-transition verification (Docker daemon unavailable so far).
- Acceptance: replay and conflicting transitions are handled safely; implementation is present, database runtime verification is pending.

## Day 3 — AI investigation (P0)
- [ ] Add provider abstraction and one tool-calling LLM integration.
- [ ] Implement four constrained read-only evidence tools.
- [ ] Persist structured reports with validated evidence references.
- [ ] Build basic investigation UI.
- Acceptance: an investigation explains a seeded discrepancy using valid application evidence.

## Day 4 — Webhooks and reconciliation (P0)
- [ ] Verify, persist, and deduplicate PayPal webhooks.
- [ ] Reconcile provider events deterministically; support status refresh.
- [ ] Add approval/rejection workflow and audit trail.
- [ ] Test duplicate and out-of-order events.
- Acceptance: events update records safely and recommendations cannot bypass backend rules.

## Day 5 — Product UI (P0/P1)
- [ ] Complete overview, payments, detail timeline, investigation, and demo screens.
- [ ] Add filtering and meaningful AG Grid integration if feasible.
- [ ] Improve loading, empty, and error states.
- Acceptance: a new user can complete the core workflow without developer assistance.

## Day 6 — Testing, documentation, deployment (P0/P1)
- [ ] Add reproducible fixtures and test critical success/failure/timeout/duplicate/discrepancy paths.
- [ ] Finalize OpenAPI, Postman, setup docs, license, and architecture docs.
- [ ] Evaluate APIMatic and Render only if stable and feasible.
- [ ] Validate clean-environment setup.
- Acceptance: core workflow passes repeatable tests and demo is reliable.

## Day 7 — Freeze and submission (P0)
- [ ] Fix high-impact bugs; review secrets, validation, authorization, and errors.
- [ ] Verify public repository, license, demo under three minutes, and official submission rules.
- [ ] Freeze features.
- Acceptance: a judge can understand, run or watch, and evaluate the product.

## Dependencies and scope
- PostgreSQL and PayPal Sandbox configuration are required for real end-to-end verification.
- AI and UI depend on stable payment/evidence APIs.
- P1 APIMatic, AG Grid, and Render are conditional. P2 analytics, Elastic, and advanced risk scoring are deferred.
