# Opportunity assessment and hard audit

## Candidate opportunities

Scores are a decision aid (1 low, 5 high), not measured user research. This assessment was first drafted 19 October 2026; delivery status below is refreshed as of 5 November 2026.

| Candidate | PayPal fit | Differentiation | Demoability | Buildability | Business case | Audit |
|---|---:|---:|---:|---:|---:|---|
| Generic PayPal AI shopping agent | 5 | 1 | 5 | 3 | 3 | Reject as lead: PayPal and many participants already cover agent shopping/action; broad scope and competitive overlap. |
| Agent spending mandate/control plane | 5 | 2 | 5 | 3 | 4 | Valuable but crowded: BOLT, Mandate, AgentLiability already occupy similar territory. |
| AI dispute-response generator | 4 | 2 | 4 | 2 | 4 | Reject as lead: PayPal Disputes APIs, Shopify AI evidence, and Chargeflow are established. Could be a later feature only with a distinct integration wedge. |
| Stablecoin/crypto payment agent | 4 | 2 | 3 | 1 | 3 | Reject for this submission: complicated flows, compliance surface and no current project integration; crypto would look ornamental without real Sandbox path. |
| Payment routing/decline prediction | 3 | 2 | 3 | 1 | 4 | Reject: no multi-processor data or labeled data; high risk of fake AI claims. |
| Evidence-first payment exception recovery | 5 | 4 | 5 | 3 | 4 | Recommend, with caveats. Existing backend covers much of the hard logic; UI/DB/Sandbox verification remain the critical delivery risks. |

## Proposed product thesis

**PayPilot — Transaction Truth for Agentic Commerce**  
Merchants and developers can see what actually happened after a PayPal payment request, when their own records disagree with PayPal or events arrive late. PayPilot reconstructs the transaction lifecycle from evidence, explains uncertainty, and prevents a retry until the outcome is verified.

### Why it belongs in this event

- PayPal is central: create/approve/capture and retrieve Orders v2 Sandbox data; verify and consume PayPal webhooks; show real provider references in the timeline.
- AI is meaningful: translate a merchant question into a bounded investigation across order/attempt/event/provider evidence; produce a structured explanation and next-step rationale. It never decides the payment state.
- Agentic angle: the scenario may originate from an agent-assisted purchase, but the unique value is accountable completion and failure handling after an agent acts.
- Trust fit: it operationalizes the public concern that AI commerce needs trusted payment and risk foundations.

### User and job

Primary user: a small merchant's developer/operator or support/finance operator without a dedicated payment reliability team. Job: “A customer says they paid, my app still shows pending, and I need to know whether it is safe to take another action.” Secondary user: engineer debugging PayPal integration behavior.

This user hypothesis needs validation. Do not imply direct PayPal customer research unless performed. A useful lightweight validation is five short interviews with merchants/developers and one coding-agent/payment-integration practitioner; ask for the last real ambiguous payment, time-to-resolution, where they checked records, and cost of duplicate capture/customer support. No leading “would you use our AI?” questions.

## Architecture change justified by the opportunity

No broad rewrite. Preserve Go/Gin, PostgreSQL, PayPal adapter, investigations package, and idempotency. The queue/detail/timeline workbench and synthetic scenario now exist. Remaining proof is the full Live API journey and real Sandbox/model behavior:

1. Exception queue: unknown/pending/review-required plus data mismatch; filters by amount, status, age.
2. Payment detail: merchant order, each capture attempt, provider IDs/status, verified webhook timeline, last refresh, investigation report.
3. One action: “Verify with PayPal” calls provider refresh; it does not create another capture while outcome is uncertain.
4. One AI action: “Investigate” runs existing constrained read-only functions. Render facts with evidence links and hypotheses in a visibly separate region.
5. Recovery state machine: retry/create another attempt only when deterministic backend rules and current provider state permit; otherwise retain manual review.
6. Scenario: one seeded timeout-after-provider-success case; clearly badge all generated/demo events as synthetic if not created in PayPal Sandbox.

Use AG Grid only if the exception queue needs sorting/filtering and there is time to implement it properly. A sponsor SDK should earn its place by improving the workflow, not appear as a logo-only checkbox.

## Business case, without invented ROI

### Who benefits

- **Merchant:** less manual stitching across logs, fewer customer-support escalations, reduced risk of duplicate capture, faster understanding of stale/inconsistent payment state.
- **Developer/operator:** shared payment timeline and consistent evidence identifiers help diagnose integration failures and support production reliability.
- **PayPal (hypothesis):** developers may build more reliable integrations and merchants may have better operational experiences when Sandbox/live API status, webhooks, and local application state are easier to reconcile. This is plausible but not validated as a PayPal objective.
- **End customer:** fewer “we don't know if you paid” exchanges and less risk of being asked to pay again.

### How to measure a business pilot

Collect a baseline first. Candidate metrics:

- Median time from exception creation to verified resolution.
- Share of exceptions resolved without manual log/provider-dashboard search.
- False paid / false failed / duplicate capture rate (target zero in deterministic tests).
- Number of support contacts per 100 payment exceptions.
- Evidence completeness and operator correction rate for AI reports.
- Sandbox task completion rate and clean-setup time for new developers.

Do **not** claim increased authorization rates, lower chargebacks, or a dollar savings number based only on a prototype. Those need production data, a control group, and appropriate PayPal/merchant permission.

## Current delivery evidence (5 November 2026)

- **Implemented and locally verified:** Go/Gin + PostgreSQL persistence, idempotency, unknown-outcome recovery guard, verified/deduplicated webhook handling, evidence-scoped investigations with deterministic fallback, embedded responsive workbench, and clearly synthetic browser scenario. `PROJECT_STATUS.md` records DB-backed integration checks and a fresh-clone Compose startup/order smoke test.
- **UI evidence:** PR #71 and issues #22/#23 record keyboard/accessibility checks, text contrast scan, and browser widths 320, 360, 390, and 640 px.
- **Still unverified:** real PayPal Sandbox checkout and webhook delivery; end-to-end Live API operator journey; live OpenAI investigation. Fake transports and synthetic UI events do not prove those external flows.
- **Submission gates:** owner license decision, final rules/eligibility confirmation, Devpost entry, public YouTube demo, and final go/no-go review remain.

## Brutal audit: why the idea can fail

1. **AI may be cosmetic.** If deterministic code alone discovers all facts and the LLM just paraphrases them, the “AI” criterion looks thin. Demonstrate why natural-language investigation and synthesis adds value; compare with a deterministic baseline and report the model's role honestly.
2. **PayPal integration is not externally verified.** Fake transport tests and local DB startup do not establish Sandbox checkout/webhook correctness. This is a top submission risk and must be resolved before claims of end-to-end provider proof.
3. **Live workflow and presentation remain unproven.** A responsive UI and functional synthetic demo now exist, but the complete Live API journey and public end-to-end recording remain open. The hackathon rejects static or non-functional demos.
4. **Significant-change eligibility.** Existing repository pre-dates the contest. The entrant must transparently show substantial post-start additions and the sponsor decides if they qualify.
5. **No actual customer discovery.** The problem is believable but not yet validated with merchants. Frame as problem hypothesis unless interviews are completed.
6. **No auth/tenant isolation.** Project README says the API is local-only. Keep demo data synthetic and isolated. Do not host a public API with real customer or payment data.
7. **Sampled competitor data is partial.** The official project gallery has not published its full list. The positioning could collide with unseen entries.
8. **Potential overreach in “agentic commerce.”** If the demo is an ordinary checkout unrelated to an agent, don't pretend it is agent-originated. Either include a small real agent decision/structured purchase request whose actual payment still uses PayPal Sandbox, or describe the product honestly as payment operations for AI-enabled commerce.
9. **Currency limitation.** Existing README says amounts are formatted using two decimal places; keep the demo on USD and label the scope.
10. **No proof PayPal would adopt it.** A hackathon prototype is not an integration partnership, product endorsement, or hiring signal. Pitch it as a useful exploration, not as a ready PayPal product.

## Go/no-go gates as of 5 November 2026

- **Local foundation:** passed with database-backed tests and clean-clone Compose startup/order smoke test, recorded in `PROJECT_STATUS.md`.
- **Workbench quality:** passed for the tested keyboard/accessibility and phone/tablet checks, recorded in PR #71 and issues #22/#23.
- **PayPal Sandbox:** not passed; buyer-approved checkout and reachable webhook delivery still require owner account access and callback setup.
- **Live workflow/model:** not passed; complete Live API browser journey and live OpenAI evidence-backed run remain unverified.
- **Submission:** not passed; license, accurate Devpost entry, public sub-three-minute YouTube demo, and final rules/eligibility review remain.

If an external gate remains blocked, keep it visible and describe it honestly. Do not treat fake transports or the synthetic browser scenario as equivalent to real Sandbox proof.
