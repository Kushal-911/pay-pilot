# Devpost submission draft

**State:** Updated 9 November 2026. **Draft only; not ready to submit.** The PayPal Sandbox checkout/capture/webhook/refresh and synthetic hosted demo have been verified, and the owner-approved MIT license is included. The owner chose not to configure live AI and approved proceeding toward submission with this limitation disclosed. No live AI investigation has been demonstrated; the official rules require an AI tool/model/platform, so this is an eligibility risk. Public video, owner review, and final rules check remain submission gates.

## Project name

PayPilot — Transaction Truth for AI-enabled Commerce

## One-line summary

An evidence-first payment operations workbench that helps merchant operators understand ambiguous PayPal payment outcomes and prevents a blind capture retry.

## Project description

When a payment capture request times out, an operator cannot conclude that the payment failed. Retrying without checking provider state risks a duplicate charge; waiting without evidence can leave the customer and merchant uncertain.

PayPilot is a local-first merchant workbench for inspecting that uncertainty. Its Go API stores merchant orders, payment attempts, investigations, and sanitized provider-event history in PostgreSQL. The PayPal Orders v2 adapter supports Sandbox order creation, capture, refresh, and webhook signature verification. A capture outcome that cannot be confirmed is kept in an unknown state, and the application requires a provider refresh before capture can be offered again. Amount/currency mismatches require review.

Operators can inspect an order-scoped timeline and request an investigation. The backend always runs deterministic record checks. When configured, an OpenAI Responses API adapter can use bounded, read-only evidence tools; its output is validated against order evidence and cannot change payment state. If model output is unavailable or invalid, the product labels and returns the deterministic fallback.

The browser-only guided scenario demonstrates an ambiguous capture and safe recovery with synthetic records. It is explicitly labeled synthetic, does not call PayPal or an AI service, and moves no money. Separately, on 9 November 2026, a USD $1 PayPal Sandbox order completed buyer approval and server-side capture; PayPilot persisted the paid state and processed the matching `PAYMENT.CAPTURE.COMPLETED` webhook. This does not demonstrate a live-money payment. The project is a prototype, has no merchant authentication or tenant isolation, and must not be used with production credentials or real merchant data.

## What is working and what remains unverified

- **Verified locally:** PostgreSQL-backed integration suite; fresh-clone Compose migrations, health checks, and order creation; workbench synthetic flow; keyboard/accessibility and narrow-screen checks. See `PROJECT_STATUS.md`.
- **Verified with PayPal Sandbox:** USD $1 buyer-approved order, server-side capture, persisted paid status/capture ID, verified/correlated/processed `PAYMENT.CAPTURE.COMPLETED` webhook, and a subsequent provider refresh that returned `COMPLETED` with matching local state.
- **Implemented but not externally verified:** live OpenAI Responses API investigation path; full Live API operator journey through the workbench and broader error-state review.
- **Required before submission:** keep the live AI limitation explicit and do not claim a model-backed feature that was not exercised. Ensure judges have a functional demo through the hosted URL or complete setup instructions; record and publish the public demo video; complete owner review and final eligibility/rules check. The owner has chosen to proceed despite the unverified AI gap and accepts the associated eligibility risk; this does not guarantee compliance or acceptance by the judges.

## PayPal and AI disclosure

- **PayPal:** PayPal Developer Platform Orders v2 Sandbox adapter for order creation, capture, provider refresh, and webhook verification. A USD $1 Sandbox buyer checkout, server-side capture, and processed capture-completed webhook were verified on 9 November 2026. No live-money transaction is part of this project.
- **AI:** an OpenAI Responses API adapter is implemented for evidence-bounded investigation. The owner chose not to configure OpenAI/API access, so no live model-backed investigation has been run. Deterministic checks remain authoritative and are the fallback; the synthetic browser scenario does not invoke an AI model. Do not present this as a demonstrated live AI integration.
- **Sponsor tools/prizes:** no sponsor-specific tool or sponsor prize is currently selected or claimed.

## Who it is for and why it matters

The intended user is a merchant developer/operator or support/finance operator investigating a payment that appears pending or uncertain. PayPilot's value hypothesis is that a correlated, evidence-linked timeline and a safe provider-state check can reduce manual record-stitching and the risk of a blind retry. This hypothesis has not been validated through merchant interviews or measured outcomes; no reduction in chargebacks, support cost, or revenue impact is claimed.

## What changed during the event period

PayPilot existed before the submission period. The public Git history records its initial PayPal checkout foundation and payment reliability/investigation work on 18 October 2026, after the 1 October period opened:

- [Initial PayPal checkout foundation (`c20e77e`)](https://github.com/Kushal-911/pay-pilot/commit/c20e77e) added the Go/Gin API, database migration base, and PayPal checkout path.
- [Payment idempotency and recovery states (`43dc1a5`)](https://github.com/Kushal-911/pay-pilot/commit/43dc1a5) added safe mutation replay and ambiguous-outcome handling.
- [Investigations and webhook reconciliation (`141bec3`)](https://github.com/Kushal-911/pay-pilot/commit/141bec3) added the evidence/investigation and verified-event foundation.

Further work merged on 5 November includes the [evidence-first workbench and capture safeguards (PR #9)](https://github.com/Kushal-911/pay-pilot/pull/9); PostgreSQL-backed event, webhook, and order-lifecycle verification (PRs [#65](https://github.com/Kushal-911/pay-pilot/pull/65), [#66](https://github.com/Kushal-911/pay-pilot/pull/66), [#67](https://github.com/Kushal-911/pay-pilot/pull/67)); liveness/readiness and configurable Compose host ports (PRs [#68](https://github.com/Kushal-911/pay-pilot/pull/68), [#69](https://github.com/Kushal-911/pay-pilot/pull/69)); timeout-to-unknown refresh-before-retry behavior (PR #70); and mobile/accessibility improvements (PR #71). The entrant must review the exact commit history, explain the changes honestly, and understand that Devpost/PayPal alone decides whether they are significant enough.

## Demo and source access

- **Public source repository:** https://github.com/Kushal-911/pay-pilot
- **Functional demo access:** follow the clean setup and run instructions in [`README.md`](README.md). Start the local Compose database/API and open the embedded workbench. Select **Demo** for the synthetic walkthrough; it is the verified path that does not require external credentials. Judges need a Docker-capable environment to run the full stack.
- **Hosted demo URL:** https://paypilot-hackathon-demo.onrender.com — synthetic browser-only simulation; no live PayPal or AI calls and no money moved. The service can sleep while idle and take about a minute to wake; see `DEPLOYMENT.md`.
- **Public YouTube video (<3 minutes):** pending; add the final public URL only after the recording shows the actual app functioning and all labels/claims have been checked.
- **Visible open-source license:** MIT, owner-approved and present at [`LICENSE`](LICENSE); third-party font and dependency notices are in [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

## Submission checklist

- [ ] Owner verifies eligibility for every entrant and confirms authority to submit.
- [x] Owner confirmed rights and selected MIT; root license and third-party notices are present.
- [x] Complete a real Sandbox checkout/capture and verify the processed webhook; the USD $1 flow and persisted paid status/capture ID were checked on 9 November 2026. Repeatability and provider refresh remain tracked in #4; no secrets or buyer details are recorded here.
- [ ] Complete and record the Live API operator journey and evidence review.
- [x] Owner chose to proceed toward submission with live AI unverified and the gap disclosed. Do not claim the hackathon AI integration requirement has been demonstrated; judges determine eligibility.
- [ ] Re-run clean-clone instructions and ensure judges can access the working project without restrictions through judging.
- [ ] Record a public YouTube demonstration under three minutes, showing PayPal and AI working end to end, using permitted media and no secrets.
- [ ] Replace draft wording with claims supported by the final verified build; confirm the Devpost countdown/rules and required fields.
- [ ] Submit through the entrant's Devpost account and record the submission confirmation.

## Official source checks

The official rules and deadline were rechecked 9 November 2026; the rules list 12 November 2026 at 12:00 p.m. Pacific Standard Time. The official rules page controls if sources conflict or change. The gallery status has not been rechecked since 5 November:

- [PayPal AI Hackathon event page](https://paypalaihackathon.devpost.com/)
- [Official rules](https://paypalaihackathon.devpost.com/rules)
- [Official FAQ](https://paypalaihackathon.devpost.com/details/faqs)
- [Event updates](https://paypalaihackathon.devpost.com/updates)
- [Project gallery](https://paypalaihackathon.devpost.com/project-gallery) — still unpublished when reviewed.
- Full dated requirements and scorecard: [`research/paypal-hackathon-dossier/hackathon-requirements.md`](research/paypal-hackathon-dossier/hackathon-requirements.md).
