# PayPilot Product Requirements Document

**Status:** Build specification for the PayPal AI Hackathon submission  
**Prepared:** 22 October 2026  
**Target:** Working, judge-runnable prototype and submission package by 11 November 2026 (submit early)  
**Product owner:** PayPilot project owner  
**Source of truth:** This document defines the target product. Existing code and documents are linked below and must not be mistaken for completed acceptance evidence.

## 1. Product summary

**PayPilot — Transaction Truth for AI-enabled Commerce** helps a merchant operator understand what happened after a PayPal payment request when the merchant's records, a delayed event, or a network timeout leave the outcome unclear. It correlates the merchant order, PayPal order/capture, payment attempts, verified webhooks, and provider refresh; presents a timestamped evidence trail; and guides a safe next step. A constrained AI investigation explains evidence and uncertainty. Deterministic backend policy and PayPal remain authoritative for payment state.

The hackathon build is a local-first merchant workbench using PayPal Sandbox. It is a prototype, not a production financial operations system, PayPal product, partnership, or endorsement.

## 2. Deadline, event constraints, and delivery window

The [official Devpost rules](https://paypalaihackathon.devpost.com/rules) and [event page](https://paypalaihackathon.devpost.com/) currently list the submission deadline as **12 November 2026, 12:00 p.m. Pacific Standard Time (PST)**. That is **13 November 2026, 1:30 a.m. India Standard Time (IST)**. Start submission on 11 November IST; verify the Devpost countdown and rules immediately before submission. The official rules page controls if event pages conflict or change.

Rules require a meaningful PayPal developer-platform Sandbox integration and an AI tool/model/platform, a working installable application consistent with the video/description, a public repository with visible open-source license, tool/model disclosure, and a public YouTube demo shorter than three minutes. Existing projects must have been significantly updated after 1 October 2026; explain the changes honestly. Judges determine whether changes qualify. Maintain access to a working demo for judging (1–15 December per current rules). Sponsor tools are optional; sponsor prizes impose separate criteria.

**Schedule guardrail:** current date leaves about three weeks. Feature work ends by 10 November. Keep 11 November for clean-run verification, video/submission and early submission. The final Devpost deadline is a hard external constraint, not a target build date.

## 3. Problem, audience, and product hypothesis

### Problem

When a capture request times out, the merchant cannot infer whether PayPal completed it. A blind retry can create duplicate-payment risk; waiting without evidence can delay fulfillment and support resolution. The operator otherwise has to stitch together merchant DB state, attempts, provider dashboard/API results, and webhook deliveries.

### Primary user

Small-merchant developer/operator or support/finance operator without a dedicated payment reliability team. Job: **“A customer says they paid, but my system says pending. Show me what happened and whether any next step is safe.”** Secondary user: engineer debugging a PayPal integration.

### Hypothesis and evidence limits

The problem and value proposition are hypotheses. The research folder reports no first-party merchant interviews, measured incident rate, ROI, or PayPal internal data. Do not claim reduced chargebacks, increased authorization, guaranteed savings, or PayPal adoption. Validate with short interviews/usability sessions if time permits and report observed measures only.

## 4. Positioning and differentiation

Position PayPilot around **post-action payment correctness and evidence-based recovery**, not generic shopping agents, an agent spending firewall, payment routing, or AI dispute writing. Publicly sampled projects include BOLT, Mandate, AgentLiability, Purser, PayRoute, and Clause-to-Cash; the Devpost gallery was unpublished when research was captured, so the landscape is partial. PayPal itself publishes Agent Toolkit and dispute capabilities. Our differentiator is the merchant-side causal timeline joining request/attempt/provider/webhook state and safe handling of ambiguous outcomes. This is a sample-based position, not a claim of an uncontested market or an internal PayPal gap.

## 5. Scope by delivery version

The repository does not have canonical V0/V1/V2 tags or a formal version history. The following labels normalize the observed product history and establish the delivery contract going forward.

### V0 — Existing foundation (implemented; runtime evidence incomplete)

- Go/Gin API, PostgreSQL schema/migrations, Docker Compose and embedded browser workbench.
- PayPal Orders v2 Sandbox adapter for OAuth, order creation, capture, refresh.
- Order/payment attempt persistence, idempotency middleware and deterministic provider request IDs.
- Unknown capture state pending provider refresh; PayPal webhook verification, deduplication, correlation and guarded transitions.
- Evidence-scoped investigations with deterministic checks, optional OpenAI Responses function calling, validated evidence IDs, and fallback.
- Synthetic demo scenario and Live API UI, OpenAPI and Postman starter collection.

This list describes code present in the current repository, not fully verified behavior. `PROJECT_STATUS.md` records missing PostgreSQL/Sandbox runtime verification and no auth/tenant isolation. Do not overstate it.

### V1 — Hackathon MVP (P0, required)

1. **Install and start:** clean clone can follow README and run app + PostgreSQL; `/health` reports DB availability.
2. **Sandbox checkout:** create merchant order, create PayPal order, approve with Sandbox buyer, server-side capture; persisted identifiers and status are visible.
3. **Exception detail:** show merchant order, amount/currency, attempts and statuses, PayPal IDs/status, event history, refresh result and investigations.
4. **Safe reconciliation:** an unknown outcome is never labeled failed or paid by timeout alone. Refresh/reconciliation occurs before any retry can be offered. Duplicate webhook IDs apply once; late pending/denied events cannot downgrade a completed capture. Mismatched amount/currency or uncorrelated events do not mark paid.
5. **Evidence-backed investigation:** user asks an order-scoped question. Deterministic checks always run. Optional AI uses only bounded, read-only evidence tools; facts/hypotheses are distinct; all cited evidence IDs resolve to records in that order. Provider/model failure returns labeled deterministic output.
6. **Demo workflow:** queue → detail/timeline → provider refresh → investigation → explain next safe step. Synthetic events are unmistakably labeled and never represented as real PayPal traffic.
7. **Submission package:** accurate setup, API/model disclosure, limitations, public repository/license, and functioning sub-three-minute video; documented significant post-1-Oct changes.

### V2 — Stretch, only after all V1 gates pass

- Hosted read-only judge demo with isolated synthetic data and no unauthenticated real merchant data.
- Better queue sorting/search and age/status filters if current UI is insufficient.
- One optional sponsor integration (AG Grid or APIMatic) only if it materially improves operator workflow and its prize criteria are met.
- Lightweight user validation and usability measures; improve labels based on findings.
- More robust migration/seed instructions and a repeatable scripted demo reset.

### Explicitly out of scope for this submission

- Production use, live merchant/customer data, multi-tenant public service, or production credentials.
- Real charge/refund/dispute automation, autonomous capture/retry, AI payment decisions, fraud prediction, chargeback outcome claims.
- Generic agent shopping, routing/decline prediction, crypto/stablecoin/Web3 flows, full disputes workflow.
- Enterprise auth/role approval system, complete audit/compliance program, advanced risk scoring, analytics/Elastic, microservices, broad provider support.
- AG Grid, APIMatic, Render, Elastic, analytics, or other sponsor additions as checkbox/logo-only work.

## 6. User journeys and functional requirements

### J1. Create and complete a Sandbox payment

Operator creates a USD order, starts PayPal checkout, follows the approval URL, and captures server-side. UI displays the local order and payment attempt identifiers and separate merchant/provider state. Browser return alone is never proof of payment. Use two-decimal currencies only in this MVP.

### J2. Resolve an ambiguous capture

Operator opens an `unknown` or `review_required` payment, sees why it is uncertain and its correlated identifiers, runs **Refresh PayPal status**, and sees the authoritative provider response and resulting deterministic reconciliation. The UI must not expose a retry/capture action while outcome is unresolved. A synthetic timeout-after-success path can demonstrate the failure mode, but all simulated details must carry a persistent “Synthetic demo” label and must not call PayPal or move money.

### J3. Inspect webhook history

Show event type, provider event ID, received time, correlation, verification/processing outcome and sanitized evidence. Deduplication and monotonic transitions are required. Do not display raw payloads or treat valid signatures as proof that the event matches the local order/amount.

### J4. Ask an investigation question

Operator submits an order-scoped question. The system presents confirmed facts linked to persisted records, hypotheses clearly labeled, unresolved questions, confidence explanation, and an enum-constrained recommendation. No AI response can mutate payments or call PayPal. Missing model credentials, timeout, invalid JSON, unknown evidence ID, invalid enum, or model outage must degrade to explicit deterministic mode without breaking payment functions.

### J5. Judge starts from a clean clone

README provides prerequisites, environment variable names (never values), Compose/migration/start steps, Sandbox buyer/business setup, test scenario, demo mode, API/model behavior, limitations, and known-good result. Public source contains no credentials, customer data, logs with sensitive content, or unlicensed material.

## 7. Product behavior and invariants

1. `paid/completed` requires a confirmed PayPal completed capture or permitted deterministic reconciliation based on verified, correlated evidence.
2. A timeout/network ambiguity produces `unknown` or review-required; it never means failed.
3. A new capture is blocked until authoritative provider state resolves an unknown outcome and policy allows it.
4. Repeated local idempotency key + identical request replays stored response; same key with different method/path/body conflicts.
5. Stable PayPal request ID identifies the same logical provider operation on retry.
6. Webhook authenticity, dedupe, correlation, amount/currency validation, and monotonic state rules are independent checks.
7. Investigation tools are fixed, read-only, order-scoped, input/result bounded and cannot issue provider calls or mutations.
8. Model content and merchant/provider text are untrusted data. Validate all report fields and evidence references against the loaded snapshot.
9. AI outage cannot break checkout, refresh, webhook processing or deterministic reporting.
10. UI must distinguish merchant state, provider state, and synthetic state; provide loading, empty, error, unavailable and mobile states.

## 8. API and data contract

Keep the existing REST API and update `openapi.yaml` and Postman collection whenever behavior changes. Existing routes (confirm exact request/response shapes against handlers/OpenAPI before modifying):

| Method and path | Purpose |
|---|---|
| `GET /health` | DB-backed health check |
| `POST /api/v1/orders` | Create internal order; idempotent |
| `GET /api/v1/orders` | List orders |
| `GET /api/v1/orders/{id}` | Order details |
| `GET /api/v1/orders/{id}/payments` | Attempt history |
| `POST /api/v1/orders/{id}/paypal-order` | Create PayPal order; idempotent |
| `POST /api/v1/payments/{id}/capture` | Server-side capture; idempotent |
| `GET /api/v1/payments/{id}/refresh` | Provider refresh/reconciliation |
| `POST /api/v1/webhooks/paypal` | Verify, persist, deduplicate and reconcile webhook |
| `POST /api/v1/investigations` | Persist order-scoped investigation; idempotent |
| `GET /api/v1/investigations/{id}` | Read investigation |
| `GET /api/v1/orders/{id}/investigations` | List order investigations |

Persist only fields needed for correlation and explanation. Keep provider event payload storage sanitized/minimal, investigations reproducible where permitted, and amounts as integer minor units. Do not expand customer/payment data retention without explicit need and privacy review.

## 9. Technical architecture and approved tools

Use the existing modular monolith. Avoid rewrite or new infrastructure unless required for a V1 acceptance failure.

| Concern | Use | Avoid / guardrail |
|---|---|---|
| Language/runtime | Go (version declared by repo/toolchain), modules | New language/service rewrite |
| HTTP | Existing Gin router and handlers | Business logic embedded in UI or route closures |
| Database | PostgreSQL via `pgx/v5` / `pgxpool`; SQL migrations in `migrations/` | ORM adoption for this deadline; direct unscoped SQL in model-facing tools |
| PayPal | Existing server-side REST adapter (`internal/paypal`), OAuth, Orders v2, Sandbox; official PayPal API docs | Browser-held PayPal secrets; production credentials; browser redirect as authority; unsupported SDK/API assumptions |
| AI | Existing OpenAI Responses API adapter behind `investigations.Model`; constrained fixed functions and schema validation | Model-selected SQL, arbitrary URLs/tools, payment mutation tools, unnecessary PII, claims without resolvable evidence |
| UI | Existing embedded HTML/CSS/vanilla JS in `cmd/api/web/` | Framework migration or large grid library unless V1 needs it and time remains |
| Local runtime | Docker Compose, PostgreSQL, `.env.example` | Commit `.env`, logs with sensitive data, production secret use |
| API contract | `openapi.yaml`, Postman collection | Undocumented behavior drift |
| Tests | Go standard testing, fake `http.RoundTripper`, DB integration when DB is available | Treating fake provider tests as Sandbox proof; silently skipping DB-gated verification |

Optional sponsors are not technical dependencies. AG Grid is allowed only for a real queue interaction need. APIMatic only if generated/validated API artifacts add concrete value. Render only after safe deployment design; current unauthenticated API must not be publicly exposed with real data. Do not add Elastic, Zapier, or other services without a V1 user need.

### Documentation sources of truth

- PayPal Orders: [Orders API integration](https://developer.paypal.com/api/rest/integration/orders-api), [Checkout Standard integration](https://developer.paypal.com/platforms/checkout/standard/integrate).
- Webhooks: [PayPal webhook integration and verification](https://developer.paypal.com/api/rest/webhooks/rest/).
- PayPal AI/agent context: [AI developer resources](https://developer.paypal.com/ai-tools/get-started), [Agent Toolkit](https://github.com/paypal/agent-toolkit), [LLM integration guide](https://developer.paypal.com/ai-tools/build-with-llm/). These are context, not substitutes for current endpoint docs.
- AI API: official OpenAI Responses API documentation; follow current official schema and tool-call semantics when touching `internal/investigations/openai.go`.
- Hackathon terms: [official rules](https://paypalaihackathon.devpost.com/rules), [event overview and submission requirements](https://paypalaihackathon.devpost.com/).
- Internal contract: `README.md`, `openapi.yaml`, `postman/PayPilot.postman_collection.json`, `DECISIONS.md`, `ROADMAP.md`, `PROJECT_STATUS.md`, and migrations.

## 10. Code and repository practices

- Preserve package boundaries: payment lifecycle in `internal/payments`, provider transport in `internal/paypal`, investigation/model integration in `internal/investigations`, UI in `cmd/api/web`, wiring in `cmd/api/main.go`.
- Keep provider calls server-side. Set bounded request timeouts, propagate request context, cap/validate response bodies, and avoid logging credentials or raw sensitive payloads.
- Use parameterized SQL, explicit transactions for coupled state changes, constraints/indexes for invariants, and forward-only numbered migrations. Do not rewrite applied migrations.
- Validate input at HTTP boundaries; use stable API errors/status codes. Keep money integer-valued in minor units and explicitly constrain MVP currencies.
- Add focused tests for every new payment transition, idempotency behavior, webhook branch, evidence validator, fallback and API contract change. Tests are required implementation practice; execute them in CI/review workflow and report what ran. Never describe unrun provider/DB paths as verified.
- Keep deterministic policy independent of model and UI. No client-side authority for paid/failed/retry state.
- Accessible labels, keyboard-operable controls, responsive layout, visible network/loading/error state, and clear synthetic/live badges.
- Keep dependencies small and pinned through Go module lock data; review license and maintenance before adding one.
- Update README/OpenAPI/Postman/roadmap/status when user-facing behavior, configuration or verification state changes.

### Conventional commits (required)

Every feature or fix must be committed using Conventional Commits: `<type>(optional-scope): imperative summary`. Examples: `feat(investigations): cite webhook evidence`, `fix(payments): block retry after unknown capture`, `docs(readme): document sandbox webhook setup`, `test(webhooks): cover duplicate delivery`. Use appropriate types such as `feat`, `fix`, `docs`, `test`, `refactor`, `build`, `ci`, `chore`. Keep each commit focused; use `!` and a `BREAKING CHANGE:` footer for breaking changes. Do not bundle unrelated features into one commit. This is the required commit rule for all future feature work.

### Documentation-to-main rule (required)

Whenever any documentation file is changed, commit that documentation change with a Conventional Commit and push it to the GitHub `main` branch (`origin/main`). Apply this to standalone documentation changes and documentation updated alongside code. Keep the documentation change focused; do not include unrelated work in its commit. Before pushing, confirm the current branch is `main` and the push can fast-forward. Never force-push. If `main` has advanced or the update cannot be pushed safely, fetch/reconcile with the latest `origin/main` and retry; report any unresolved blocker rather than leaving the documentation change silently unpushed. This rule applies to future documentation edits and to this PRD change.

## 11. Acceptance criteria and judge proof

### P0 functional acceptance

- Clean setup works from the documented clone path; `/health` reflects DB status.
- A real PayPal Sandbox order is created, buyer-approved, and captured server-side; resulting provider IDs/status persist and appear in UI.
- Unknown capture remains unresolved until provider refresh/event evidence establishes outcome; no duplicate capture path is offered prematurely.
- Duplicate verified webhook delivery is applied once; out-of-order event cannot downgrade completed state; mismatch is review-required, not paid.
- Investigation uses actual persisted evidence; every displayed citation resolves to the displayed order; invalid model output falls back deterministically.
- Demo clearly distinguishes real Sandbox interaction from synthetic event injection.
- Under-three-minute video shows the application functioning and makes the problem, PayPal role, AI role, and recovery outcome understandable.
- Public repo includes source, setup instructions and visible license; submission identifies AI tool/model and role; post-event-start changes are documented.

### Judge proof matrix

| Criterion | Proof artifact |
|---|---|
| Technological implementation | Sandbox walkthrough, persisted provider IDs, verified webhook/investigation evidence, test results |
| Design | Coherent queue → detail → timeline → investigation → safe next step, with responsive/error states |
| Potential impact | Named operator job, explicit hypothesis, honest validation/metrics (no invented ROI) |
| Innovation | Contrast with sampled agent-control/dispute/routing projects; emphasize cross-record post-action correctness |
| Presentation | Clear <3 minute public YouTube walkthrough, accurate labels and claims |

## 12. Delivery plan and gates

Research supplied these gates; dates below remain the working plan, not event requirements:

| Gate | Target | Required exit |
|---|---:|---|
| A: prove the foundation | 22 Oct | Clean DB/migrations and real Sandbox create/approve/capture; document any external blocker honestly |
| B: end-to-end workbench | 27 Oct | Real backend queue/detail/refresh/investigation journey |
| C: prove safety | 31 Oct | Unknown outcome never recaptures; duplicate/out-of-order webhook and evidence correctness verified |
| D: clean install | 4 Nov | Independent clean-run, demo reset, truthful synthetic labels and setup guide |
| E: submission ready | 7 Nov | Video and text reviewed against runtime; significant changes documented |
| Feature freeze | 10 Nov | Only correctness and submission blockers may change |
| Submit early | 11 Nov IST | Public repository/license, video and complete Devpost submission submitted; confirm receipt |
| Hard deadline | 13 Nov 1:30 a.m. IST | Devpost currently lists Nov 12 noon PST; re-check countdown and timezone |

If gates slip, cut stretch scope and polish before weakening payment invariants. A minimum credible fallback is a real successful Sandbox checkout plus one repeatable, explicitly synthetic ambiguous-outcome investigation with safe refresh behavior; disclose any unverified webhook/DB functionality.

## 13. Metrics and validation

For prototype usability, measure time to find provider state, whether operator distinguishes fact from hypothesis, understanding of why retry is blocked/allowed, wrong-action rate, evidence click-through/correction rate, and setup completion. For a pilot, baseline median/P90 time-to-resolution, manual lookup steps, support contacts per exception, evidence completeness, model fallback rate/latency, and false state/duplicate action rate (zero in deterministic tests). Do not claim commercial impact until measured with an appropriate baseline and permission.

## 14. Risks and open constraints

- Runtime dependencies/credentials may block DB and Sandbox proof; fake transport tests are not equivalent.
- Existing API has no merchant authentication or tenant isolation. Keep deployment local or synthetic/read-only and do not use real customer data.
- Current currency conversion assumes two decimal places; keep MVP to USD.
- Existing project predates hackathon; preserve commits/release notes/screenshots demonstrating meaningful post-1-Oct work. Eligibility remains sponsor/administrator decision.
- Public Devpost gallery research is incomplete; differentiation may overlap unseen entries.
- AI can appear cosmetic unless the demo shows valuable natural-language evidence synthesis; compare claims against deterministic output.
- Avoid claims about PayPal internal systems, adoption, or partnership. Research supports only public-facing context and explicitly marks strategic fit as inference.

## 15. Research and artifact traceability

Research provenance and all conclusions are indexed here. Each source is an input; current acceptance status is defined by the code plus fresh verification, not by an older status note.

| Research / existing artifact | PRD use |
|---|---|
| `research/paypal-hackathon-dossier/README.md` | Research map and reading order |
| `research/paypal-hackathon-dossier/hackathon-requirements.md` | Deadline, rules, eligibility, requirements, scoring, submission gates |
| `research/paypal-hackathon-dossier/sources.md` | Primary source ledger, external product/API/competitor sources, limitations |
| `research/paypal-hackathon-dossier/project-landscape.md` | Competitive overlap and positioning |
| `research/paypal-hackathon-dossier/opportunity-analysis.md` | Product thesis, users, scope, business hypothesis, go/no-go gates and risks |
| `research/paypal-hackathon-dossier/paypal-direction.md` | Public PayPal product context and fact/inference boundary |
| `research/paypal-hackathon-dossier/pitch-and-demo.md` | Narrative and sub-three-minute proof sequence |
| `research/paypal-hackathon-dossier/people-and-outreach.md` | Optional feedback channels and outreach guardrails; no outreach is implied by this PRD |
| `research/paypilot-technical-playbook/README.md` | Architecture/playbook navigation and boundary |
| `research/paypilot-technical-playbook/architecture.md` | Current and target architecture, state authority and production gaps |
| `research/paypilot-technical-playbook/implementation-guide.md` | Tool choices, sequence, learning path, verification checklist |
| `research/paypilot-technical-playbook/payment-reliability.md` | Ambiguous outcomes, idempotency, webhook and reconciliation invariants |
| `research/paypilot-technical-playbook/ai-evidence.md` | Read-only AI boundary, evidence types, validation and fallback |
| `research/paypilot-technical-playbook/business-case.md` | Hypothesis, value chain, claims to avoid, discovery and pilot metrics |
| `README.md`, `PROJECT_STATUS.md`, `HANDOFF.md` | Existing product behavior, implementation audit and known blockers |
| `ROADMAP.md`, `DECISIONS.md`, `DEMO_SCRIPT.md` | Milestones, architecture decisions, demo claims and script |
| `openapi.yaml`, `postman/PayPilot.postman_collection.json` | Existing API surface and API contract artifacts |
| `cmd/api/main.go`, `internal/payments/`, `internal/paypal/`, `internal/investigations/`, `migrations/`, `cmd/api/web/` | Current implementation to preserve, inspect and bring to acceptance |

External reference list: [Devpost rules](https://paypalaihackathon.devpost.com/rules); [event page](https://paypalaihackathon.devpost.com/); PayPal [Orders API](https://developer.paypal.com/api/rest/integration/orders-api), [Checkout](https://developer.paypal.com/platforms/checkout/standard/integrate), [webhooks](https://developer.paypal.com/api/rest/webhooks/rest/), [AI resources](https://developer.paypal.com/ai-tools/get-started), [Agent Toolkit](https://github.com/paypal/agent-toolkit), [LLM guide](https://developer.paypal.com/ai-tools/build-with-llm/). The full discovery and competitor bibliography remains in `research/paypal-hackathon-dossier/sources.md`.
