# PayPilot delivery roadmap and issue register

**Updated:** 5 November 2026  
**Repository:** [Kushal-911/pay-pilot](https://github.com/Kushal-911/pay-pilot)  
**Product requirements:** [PRD.md](PRD.md)  
**Hackathon deadline currently published:** 12 November 2026, 12:00 p.m. PST / 13 November 2026, 1:30 a.m. IST. Recheck the [official Devpost rules](https://paypalaihackathon.devpost.com/rules) and countdown before submission.

This is the durable index for the complete planned lifecycle. V0 is the foundation already present; V1 is the bounded hackathon deliverable; V2 is optional post-V1 demo polish; V3 is the longer production-readiness target. **Creating an issue does not mean its feature is implemented or accepted.** Close issues only after recording their acceptance evidence.

## Current status

- **V0 code exists; runtime proof is incomplete.** The repository contains the Go/Gin API, PostgreSQL migrations, PayPal Sandbox adapter, idempotency, safe unknown-capture behavior, webhook verification/deduplication, evidence-bounded investigations, embedded workbench, synthetic demo, and API artifacts. See [PROJECT_STATUS.md](PROJECT_STATUS.md) for the verification gaps.
- **UI implementation is substantially present.** Demo and Live API workflows, order queue/detail, payment timeline, investigation view, and new-payment flow exist. That does not make the end-to-end product complete: PostgreSQL/Sandbox/runtime evidence is still required.
- **Feature work in review:** [PR #9](https://github.com/Kushal-911/pay-pilot/pull/9) is open and unmerged. It covers safer payment actions, provider event timeline, truthful Live API navigation/metrics, claim-level AI evidence, capture amount/currency verification, and bounded PayPal requests/responses. Branch commit `3b8a7bb` passes `go test ./...`, `go vet ./...`, and `git diff --check`. Issues #1, #2, #5, #7, #16, and #17 retain integration/runtime acceptance.
- **Unfinished V1 gates:** clean database-backed run; real Sandbox buyer checkout and webhook path; verified safe recovery; API contract synchronization; clean-clone run; owner-confirmed license; submission description; public demo video; final requirements check. Docker cannot connect to its local daemon, database URL variables and PayPal Sandbox credentials are absent, so database and provider runtime gates are blocked in this environment. The duplicate route in PR #9 has been removed in commit `5aa515e`; startup verification remains open in #62.
- **Backlog:** 61 GitHub issues are open in the repo, including this index's original issues #1–#8 and 53 newly created work/epic issues #10–#62. Pull request #9 occupies the number between them.

## Phases and completion criteria

| Phase | Tracker | Definition of done |
|---|---|---|
| V0 — Foundation | [#58](https://github.com/Kushal-911/pay-pilot/issues/58) | Accurately document the baseline and its tested versus unverified behavior. V0 is not proof of a working external integration. |
| V1 — Hackathon MVP | [#59](https://github.com/Kushal-911/pay-pilot/issues/59) | Fresh-clone run, PostgreSQL-backed lifecycle, real Sandbox checkout and webhook evidence, safe ambiguous-outcome recovery, evidence-grounded investigation or explicit deterministic fallback, consistent API docs, owner-approved license, accurate submission and public video. |
| V2 — Judge Demo | [#60](https://github.com/Kushal-911/pay-pilot/issues/60) | Optional isolated synthetic hosted demo, reset and observed usability improvements. Start only after V1; V2 must not delay V1. |
| V3 — Production Readiness | [#61](https://github.com/Kushal-911/pay-pilot/issues/61) | Security, tenant isolation, privacy, reliability, operational and deployment gates evidenced before real customer data or a production claim. This is long-horizon scope, not a hackathon promise. |

**Native GitHub milestones are not created.** The connected GitHub issue controls can create/update issues and apply existing labels, but do not expose milestone creation, custom label creation, or Projects v2 management. The repo Milestones page also showed a signed-out state. The four issues above are clearly titled milestone epics and all work issues carry phase tags. Once milestone administration is available, create exact milestone names `V0 — Foundation`, `V1 — Hackathon MVP`, `V2 — Judge Demo`, and `V3 — Production Readiness`, then bulk-assign issues by title tag. No issue is falsely presented as having a native milestone assignment.

## Priority and sequencing

- **V1 P0 — required before submission:** #1–#4, #6 and #10–#18, #25–#29, #62. #1/#2 have implementation in PR #9; #3/#4 and runtime evidence remain open. #62 tracks startup verification for the duplicate route removed from PR #9 in commit `5aa515e`. Owner-only license and submission actions remain explicit blockers.
- **V1 P1 — polish/evidence:** #5, #7, #19–#24 and #30–#31. Confirm the live journey, model failure behavior, accessibility/responsive behavior, synthetic demo repeatability, user observations where feasible, and current event rules.
- **V2 — after V1 only:** #32–#36. Hosted demo must be synthetic/read-only and isolated.
- **V3 — production path:** #37–#57. These tasks address the current local prototype's lack of authentication/tenant isolation and add payment operations, privacy, reliability and security controls.
- All issues include a concrete goal/scope, acceptance criteria, dependencies or blockers, and pointers to the source-of-truth documents/research where relevant. See each issue body for details.

## Label and naming conventions

New work issues use the repository's existing GitHub type labels: `enhancement`, `bug`, or `documentation`. Version, priority, and work area are encoded in issue titles, such as `[V1][P0][payments]`. No custom label names are claimed or fabricated. Priority/phase/area labels can be added later if repo label administration is enabled.

Issues are not assigned to a person because no confirmed assignee mapping was provided. Owner-only decisions (license, credentials, hosted account and Devpost/video publishing) are called out as dependencies rather than guessed.

## V1 — Hackathon MVP and submission

- [#1 — [V1 P0] Make Live API payment actions distinguish capture from provider refresh](https://github.com/Kushal-911/pay-pilot/issues/1)
- [#2 — [V1 P0] Expose verified PayPal event history to the payment detail timeline](https://github.com/Kushal-911/pay-pilot/issues/2)
- [#3 — [V1 P0] Prove PostgreSQL-backed lifecycle and webhook behavior](https://github.com/Kushal-911/pay-pilot/issues/3)
- [#4 — [V1 P0 BLOCKED] Complete a real PayPal Sandbox checkout and refresh](https://github.com/Kushal-911/pay-pilot/issues/4)
- [#5 — [V1 P1] Replace placeholder workbench navigation and static overview figures](https://github.com/Kushal-911/pay-pilot/issues/5)
- [#6 — [V1 P0] Finish clean-clone documentation and hackathon submission package](https://github.com/Kushal-911/pay-pilot/issues/6)
- [#7 — [V1 P1] Demonstrate and validate the AI investigation's operator value](https://github.com/Kushal-911/pay-pilot/issues/7)
- [#10 — [V1][P0][runtime] Start the full stack from a clean checkout](https://github.com/Kushal-911/pay-pilot/issues/10)
- [#11 — [V1][P0][payments] Exercise idempotency replay and conflict against PostgreSQL](https://github.com/Kushal-911/pay-pilot/issues/11)
- [#12 — [V1][P0][payments] Verify order and payment lifecycle in PostgreSQL](https://github.com/Kushal-911/pay-pilot/issues/12)
- [#13 — [V1][P0][webhooks] Verify webhook deduplication and monotonic transitions in PostgreSQL](https://github.com/Kushal-911/pay-pilot/issues/13)
- [#14 — [V1][P0][paypal] Complete the Sandbox checkout using a test buyer](https://github.com/Kushal-911/pay-pilot/issues/14)
- [#15 — [V1][P0][paypal] Verify webhook delivery through a reachable Sandbox callback](https://github.com/Kushal-911/pay-pilot/issues/15)
- [#16 — [V1][P0][payments] Enforce amount and currency matching before paid state](https://github.com/Kushal-911/pay-pilot/issues/16)
- [#17 — [V1][P0][payments] Cover provider timeout, malformed response, and retry-safe recovery](https://github.com/Kushal-911/pay-pilot/issues/17)
- [#18 — [V1][P0][api] Sync OpenAPI and Postman with the shipped API](https://github.com/Kushal-911/pay-pilot/issues/18)
- [#19 — [V1][P1][ui] Validate the complete operator journey in Live API mode](https://github.com/Kushal-911/pay-pilot/issues/19)
- [#20 — [V1][P1][investigation] Test invalid and unavailable model responses](https://github.com/Kushal-911/pay-pilot/issues/20)
- [#21 — [V1][P1][investigation] Verify a live AI investigation against persisted Sandbox evidence](https://github.com/Kushal-911/pay-pilot/issues/21)
- [#22 — [V1][P1][quality] Add accessibility and keyboard checks to the workbench](https://github.com/Kushal-911/pay-pilot/issues/22)
- [#23 — [V1][P1][quality] Verify narrow-screen layout for the recovery workflow](https://github.com/Kushal-911/pay-pilot/issues/23)
- [#24 — [V1][P1][ui] Make the synthetic failure scenario repeatable and unmistakable](https://github.com/Kushal-911/pay-pilot/issues/24)
- [#25 — [V1][P0][docs] Prove the README from a fresh clone](https://github.com/Kushal-911/pay-pilot/issues/25)
- [#26 — [V1][P0][legal] Confirm project ownership and select the repository license](https://github.com/Kushal-911/pay-pilot/issues/26)
- [#27 — [V1][P0][submission] Prepare an accurate Devpost submission draft](https://github.com/Kushal-911/pay-pilot/issues/27)
- [#28 — [V1][P0][release] Run the submission go/no-go checklist and freeze features](https://github.com/Kushal-911/pay-pilot/issues/28)
- [#29 — [V1][P0][submission] Record and review the public demo under three minutes](https://github.com/Kushal-911/pay-pilot/issues/29)
- [#30 — [V1][P1][research] Run lightweight operator usability sessions](https://github.com/Kushal-911/pay-pilot/issues/30)
- [#31 — [V1][P1][research] Recheck hackathon rules, gallery, and sponsor criteria before submission](https://github.com/Kushal-911/pay-pilot/issues/31)
- [#62 — [V1][P0][bug] Remove duplicate provider-event route registration](https://github.com/Kushal-911/pay-pilot/issues/62)

## V2 — Hosted judge demo and validated polish

- [#32 — [V2][P0][architecture] Specify a safe hosted judge demo](https://github.com/Kushal-911/pay-pilot/issues/32)
- [#33 — [V2][P1][deploy] Deploy the isolated synthetic judge demo](https://github.com/Kushal-911/pay-pilot/issues/33)
- [#34 — [V2][P1][demo] Add a one-command deterministic demo reset](https://github.com/Kushal-911/pay-pilot/issues/34)
- [#35 — [V2][P2][ux] Improve queue search and filters from observed operator needs](https://github.com/Kushal-911/pay-pilot/issues/35)
- [#36 — [V2][P2][sponsor] Evaluate optional sponsor integrations against a user need](https://github.com/Kushal-911/pay-pilot/issues/36)

## V3 — Production readiness

- [#37 — [V3][P0][security] Define authentication and session boundary for merchant access](https://github.com/Kushal-911/pay-pilot/issues/37)
- [#38 — [V3][P0][security] Enforce tenant isolation in every data access path](https://github.com/Kushal-911/pay-pilot/issues/38)
- [#39 — [V3][P0][security] Add authorization policy for payment operations](https://github.com/Kushal-911/pay-pilot/issues/39)
- [#40 — [V3][P0][security] Add secret rotation and deployment secret management](https://github.com/Kushal-911/pay-pilot/issues/40)
- [#41 — [V3][P0][security] Add abuse controls to public API endpoints](https://github.com/Kushal-911/pay-pilot/issues/41)
- [#42 — [V3][P1][audit] Persist an append-only operator action audit trail](https://github.com/Kushal-911/pay-pilot/issues/42)
- [#43 — [V3][P0][payments] Support currency-specific minor-unit precision](https://github.com/Kushal-911/pay-pilot/issues/43)
- [#44 — [V3][P1][webhooks] Add durable webhook retry and recovery workflow](https://github.com/Kushal-911/pay-pilot/issues/44)
- [#45 — [V3][P1][reconciliation] Schedule safe reconciliation for unresolved outcomes](https://github.com/Kushal-911/pay-pilot/issues/45)
- [#46 — [V3][P1][operations] Add structured logs, metrics, and correlation tracing](https://github.com/Kushal-911/pay-pilot/issues/46)
- [#47 — [V3][P0][privacy] Define data minimization, retention, deletion, and export policy](https://github.com/Kushal-911/pay-pilot/issues/47)
- [#48 — [V3][P1][reliability] Add backup, restore, and disaster recovery procedures](https://github.com/Kushal-911/pay-pilot/issues/48)
- [#49 — [V3][P1][reliability] Define SLOs and incident runbooks for payment recovery](https://github.com/Kushal-911/pay-pilot/issues/49)
- [#50 — [V3][P0][quality] Build a security and threat-model test suite](https://github.com/Kushal-911/pay-pilot/issues/50)
- [#51 — [V3][P1][quality] Establish load and resilience targets for critical workflows](https://github.com/Kushal-911/pay-pilot/issues/51)
- [#52 — [V3][P0][delivery] Add CI gates for build, tests, migrations, and dependency hygiene](https://github.com/Kushal-911/pay-pilot/issues/52)
- [#53 — [V3][P1][delivery] Define forward-only schema migration and rollback procedures](https://github.com/Kushal-911/pay-pilot/issues/53)
- [#54 — [V3][P1][architecture] Review PayPal API versioning, idempotency, and deprecation exposure](https://github.com/Kushal-911/pay-pilot/issues/54)
- [#55 — [V3][P1][operations] Add authorized dead-letter inspection and replay controls](https://github.com/Kushal-911/pay-pilot/issues/55)
- [#56 — [V3][P2][product] Decide whether disputes belong in a separately scoped product phase](https://github.com/Kushal-911/pay-pilot/issues/56)
- [#57 — [V3][P2][product] Establish a production-readiness review gate](https://github.com/Kushal-911/pay-pilot/issues/57)

## Milestone tracker epics

- [#58 — [MILESTONE V0][Foundation] Baseline implementation and verified local prototype](https://github.com/Kushal-911/pay-pilot/issues/58)
- [#59 — [MILESTONE V1][Hackathon MVP] Complete and submit a working PayPilot prototype](https://github.com/Kushal-911/pay-pilot/issues/59)
- [#60 — [MILESTONE V2][Judge Demo] Isolated hosted synthetic demonstration and usability improvements](https://github.com/Kushal-911/pay-pilot/issues/60)
- [#61 — [MILESTONE V3][Production Readiness] Secure, operable multi-tenant payment recovery service](https://github.com/Kushal-911/pay-pilot/issues/61)

## Product completion

**Hackathon MVP is complete** only after all V1 P0 acceptance evidence exists: an independent clean clone starts the app; migrations and database-gated integration checks run; a Sandbox buyer completes checkout and server-side capture; webhook verification/correlation and unknown-outcome refresh are demonstrated; the workbench displays persisted evidence and safe next actions; AI output is order-scoped and evidence-linked or deterministic fallback is explicit; OpenAPI/Postman and setup docs match; the owner approves the license; and the public video is under three minutes and accurately discloses synthetic behavior and the AI role.

**The full project target** is the V3 production-readiness gate, not merely a polished UI. V3 completion requires closed security, privacy, operational, resilience and deployment issues with owner review. It does not automatically imply regulatory certification, PayPal endorsement, guaranteed financial impact, or suitability for production funds.

## Repository and code practices

- Create or identify an issue before a feature. Link the pull request to its issue; update the issue with implementation, verification and remaining blockers. Close only when acceptance criteria pass.
- Use Conventional Commits: `feat(payments): reconcile unknown capture outcomes`, `fix(webhooks): reject mismatched amounts`, `docs(api): sync event history contract`, `test(payments): cover idempotency conflict`.
- Use focused branches and reviews. Do not force-push.
- Any documentation change must use a Conventional Commit and be pushed directly to `main`, per [AGENTS.md](AGENTS.md). Keep the documentation commit focused. Before pushing, verify the branch is `main` and the push is a fast-forward.
- Follow [PRD.md §10](PRD.md#10-code-and-repository-practices): preserve Go package boundaries; use parameterized SQL and explicit transactions; add forward-only numbered migrations; represent money in integer minor units; bound provider/model calls; validate order-scoped evidence; sanitize logs/payloads; test duplicate, error, and unknown-state behavior; never commit credentials.
- Keep PayPal and AI calls server-side. Label synthetic data persistently. Use Sandbox credentials only; do not connect real merchant/customer data to the current unauthenticated prototype.

## Source-of-truth and research index

- Product behavior, versions, scope, architecture and code practices: [PRD.md](PRD.md)
- Implementation status and blockers: [PROJECT_STATUS.md](PROJECT_STATUS.md)
- Existing project plan: [ROADMAP.md](ROADMAP.md)
- Install/run guide: [README.md](README.md)
- Demo walkthrough: [DEMO_SCRIPT.md](DEMO_SCRIPT.md)
- Decision log: [DECISIONS.md](DECISIONS.md)
- Technical build guide: [research/paypilot-technical-playbook/README.md](research/paypilot-technical-playbook/README.md)
- Hackathon requirements/entry audit: [research/paypal-hackathon-dossier/hackathon-requirements.md](research/paypal-hackathon-dossier/hackathon-requirements.md)
- Research source ledger and limitations: [research/paypal-hackathon-dossier/sources.md](research/paypal-hackathon-dossier/sources.md)
- Official PayPal Orders/Checkout, webhook, PayPal AI, OpenAI Responses API and Devpost rules links are collected in [PRD.md §9](PRD.md#9-technical-architecture-and-approved-tools) and the source ledger. Verify current provider docs before changing an integration.
