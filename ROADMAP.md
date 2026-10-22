# PayPilot Roadmap

Status legend: `[x]` verified complete, `[~]` implemented or partially complete but acceptance not verified, `[ ]` not started or blocked.

## Day 1 — Foundation and PayPal checkout (P0)
- [~] Initialize Go/Gin backend and PostgreSQL persistence (code exists; PostgreSQL runtime unverified).
- [~] Add Docker Compose, environment example, health endpoint, and initial SQL migration (config validates; daemon unavailable).
- [~] Implement internal order creation and retrieval (not exercised against PostgreSQL).
- [~] Implement PayPal order creation, approval URL response, and server-side capture endpoint (Sandbox flow not run).
- [~] Start Postman collection (starter requests exist; live API not run).
- [ ] Exercise the complete Sandbox checkout with configured credentials and buyer approval.
- Acceptance: a Sandbox payment completes and its verified result is persisted and visible. External credentials required.

## Day 2 — Payment lifecycle and reliability (P0)
- [~] Add payment attempts, idempotency, state-transition rules, and payment history/detail (code exists; DB integration pending).
- [~] Model unknown provider outcomes and safe refresh before consequential retries (unit adapter path tested; DB transitions unverified).
- [~] Add focused adapter tests and a PostgreSQL-gated idempotency integration test (unit tests pass; DB test skipped without `TEST_DATABASE_URL`).
- [ ] Run PostgreSQL integration and state-transition verification (Docker daemon unavailable so far).
- Acceptance: replay and conflicting transitions are handled safely; implementation is present, database runtime verification is pending.

## Day 3 — AI investigation (P0)
- [~] Add provider abstraction and optional OpenAI Responses tool-calling integration (fake API tested; no live credentials).
- [~] Implement four constrained read-only evidence tools.
- [~] Persist structured reports with validated evidence references (migration/handler implemented; DB path unverified).
- [~] Build investigation UI with evidence-linked results; browser and persisted API flow pending.
- Acceptance: backend report and evidence workflow implemented; full persisted discrepancy demo and UI remain unverified.

## Day 4 — Webhooks and reconciliation (P0)
- [~] Verify, persist, and deduplicate PayPal webhooks (postback verification implemented; PayPal/PostgreSQL runtime pending).
- [~] Reconcile order/attempt/provider event records deterministically; fake out-of-order tests pass; DB runtime pending.
- [ ] Add approval/rejection workflow and audit trail.
- [~] Test out-of-order event decisions with pure tests; duplicate-event DB behavior still needs integration verification.
- Acceptance: events update records safely and recommendations cannot bypass backend rules.

## Day 5 — Product UI (P0/P1)
- [~] Complete responsive overview, payments, detail timeline, investigation, and guided demo screens in an embedded workbench.
- [~] Add payment search/filter; AG Grid deferred because the current queue is small and dependency-free.
- [~] Add loading, empty, error, API-unavailable, and mobile states.
- Acceptance: synthetic scenario is walkthrough-ready; Live API usability still needs runtime validation.

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
