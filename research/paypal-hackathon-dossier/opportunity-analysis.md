# Opportunity assessment and hard audit

## Candidate opportunities

Scores are a decision aid (1 low, 5 high), not measured user research. “Buildability” assumes the existing PayPilot codebase and roughly three weeks until deadline.

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

No broad rewrite. Preserve Go/Gin, PostgreSQL, PayPal adapter, investigations package, and idempotency. Add only the user-facing workflow and real scenario observability:

1. Exception queue: unknown/pending/review-required plus data mismatch; filters by amount, status, age.
2. Payment detail: merchant order, each capture attempt, provider IDs/status, verified webhook timeline, last refresh, investigation report.
3. One action: “Reconcile with PayPal” calls existing refresh; it does not create another capture.
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

## Brutal audit: why the idea can fail

1. **AI may be cosmetic.** If deterministic code alone discovers all facts and the LLM just paraphrases them, the “AI” criterion looks thin. Demonstrate why natural-language investigation and synthesis adds value; compare with a deterministic baseline and report the model's role honestly.
2. **PayPal integration may be unverified.** Fake transport tests do not establish sandbox correctness. The existing audit reports no database runtime and no PayPal sandbox test. This is the top risk and must be resolved before polish.
3. **UI absence threatens design and demo.** The backend has no merchant UI. The hackathon explicitly rejects static prototypes in place of a functioning app.
4. **Significant-change eligibility.** Existing repository pre-dates the contest. The entrant must transparently show substantial post-start additions and the sponsor decides if they qualify.
5. **No actual customer discovery.** The problem is believable but not yet validated with merchants. Frame as problem hypothesis unless interviews are completed.
6. **No auth/tenant isolation.** Project README says the API is local-only. Keep demo data synthetic and isolated. Do not host a public API with real customer or payment data.
7. **Sampled competitor data is partial.** The official project gallery has not published its full list. The positioning could collide with unseen entries.
8. **Potential overreach in “agentic commerce.”** If the demo is an ordinary checkout unrelated to an agent, don't pretend it is agent-originated. Either include a small real agent decision/structured purchase request whose actual payment still uses PayPal Sandbox, or describe the product honestly as payment operations for AI-enabled commerce.
9. **Currency limitation.** Existing README says amounts are formatted using two decimal places; keep the demo on USD and label the scope.
10. **No proof PayPal would adopt it.** A hackathon prototype is not an integration partnership, product endorsement, or hiring signal. Pitch it as a useful exploration, not as a ready PayPal product.

## Go/no-go gates

- **Gate A (22 Oct):** clean local DB startup/migrations + Sandbox create, approve, capture from a fresh run. If blocked, use verified publicly documented Sandbox test flow and keep the demo's external limitation explicit; don't claim live integration.
- **Gate B (27 Oct):** UI connects to the real backend and completes order -> exception -> refresh -> investigation with traceable records.
- **Gate C (31 Oct):** timeout/unknown scenario proven not to issue a second capture; duplicate/out-of-order webhooks do not downgrade paid state; evidence IDs resolve to actual records.
- **Gate D (4 Nov):** independent user can install and run from README; one demo scenario is clear and all synthetic records labeled.
- **Gate E (7 Nov):** 3-minute video proves function, no unverified claims; submission explains post-1-Oct progress.

If a gate fails, cut features, not correctness. The minimal acceptable prototype is a single successful Sandbox checkout plus one reproducible ambiguous-outcome investigation and safe refresh behavior.
