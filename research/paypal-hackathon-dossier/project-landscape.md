# Project and competition landscape

**Snapshot date:** 19 October 2026. The official [gallery currently says it has not been published](https://paypalaihackathon.devpost.com/project-gallery). Publicly indexed individual entries nevertheless provide an early sample. This table is deliberately selective, not a claim to cover every entrant.

## Current event entries found

| Project | Public proposition | Overlap / implication for PayPilot |
|---|---|---|
| [BOLT Agentic Payments Control Plane](https://devpost.com/software/bolt-agentic-payments-control-plane) | Allow/approval/deny policies, spend limits, approvals, audit for AI-originated PayPal payment/refund actions | Strong overlap with generic “AI must not autonomously move money” framing. Don't pitch PayPilot as just a payment-agent policy firewall. |
| [Mandate](https://devpost.com/software/mandate-n7h1x5) | Agent shopping constrained by a user mandate; experiments showing prompt-only controls and overpayment | Direct overlap with budget/intent controls. It emphasizes policy against overpaying, not post-timeout provider-state reconstruction. |
| [AgentLiability](https://devpost.com/software/agentliability) | AI policy layer, PayPal transactions, risk assessment, dispute logging | Further crowds the guardrail/liability framing. Differentiate using causal payment evidence and verified recovery decisions. |
| [Purser](https://devpost.com/software/purser) | Group-event treasurer: ledger, invoices/orders, approvals, refunds, audit, PayPal Agent Toolkit, Gemini, AG Grid | A sophisticated complete workflow; shows that the bar is not “API plus chatbot.” Avoid group expense or simple approval-ledger ideas. |
| [PayRoute](https://devpost.com/software/payroute-ai-powered-smart-payment-routing) | AI-supported smart payment routing | Another common payments-AI pattern. PayPilot is not routing or underwriting. |
| [Clause-to-Cash](https://devpost.com/software/clause-to-cash-contracts-that-collect) | Extract payment terms from contracts, cite source pages, then help collect | Demonstrates strong citation/provenance UX. PayPilot must make citations clickable to actual transaction rows/events, not merely quote documents. |

### Earlier adjacent projects and PayPal demonstrations

- **CrawlerPal**, the first-place project at PayPal Dev Days 2025, used PayPal APIs to monetize AI bot/API access via an HTTP 402/L402-style paywall. PayPal's [winner recap](https://developer.paypal.com/community/blog/winners-dev-days-hackathon/) describes it. This signals interest in agent-to-agent commerce and meaningful PayPal use, not just a generic consumer assistant.
- **PayPal Agent Toolkit** itself offers a broad action layer. The [official repo](https://github.com/paypal/agent-toolkit) and [release note](https://developer.paypal.com/community/blog/paypal-agentic-ai-toolkit/) make shallow action wrappers crowded territory.
- PayPal Developer's LinkedIn account has amplified a winning voice agent that rebooks after a flight cancellation, combining user-facing service automation with PayPal ([public post](https://www.linkedin.com/posts/paypaldev_we-just-had-to-share-with-the-community-this-activity-748815031488165-t5B-)). The exact URL may redirect; validate before quoting in a presentation.

## Commercial alternatives and incumbents

1. **PayPal-native disputes.** PayPal already offers a [Disputes API](https://developer.paypal.com/api/disputes/) and reason-specific [evidence guidance](https://developer.paypal.com/disputes/reasons-evidence/). The Agent Toolkit includes dispute list/get and some actions. A simple “AI writes a dispute response” feature is therefore a poor lead idea.
2. **Shopify.** Shopify collects/submits chargeback evidence in Admin and includes optional AI-generated insights for “product not received” disputes ([help documentation](https://help.shopify.com/en/manual/payments/chargebacks/chargebacks-in-admin)).
3. **Chargeflow.** Chargeflow markets end-to-end ingestion, enrichment, evidence creation, and submission for disputes across processors ([API automation docs](https://docs.chargeflow.io/docs/merchants/automation)). This is a mature direct competitor to broad AI chargeback automation.
4. **Riskified / Visa.** Riskified markets AI-driven dispute resolution for Shopify; Visa has AI-enhanced dispute products and pre-dispute workflows ([Visa announcement](https://usa.visa.com/about-visa/newsroom/press-releases.releaseId.22261.html)). This makes “AI dispute automation” an occupied category.

## Pattern analysis

Current public entries frequently present an agent that proposes or executes a new transaction, then add a deterministic policy guard. That is relevant and technically valuable, but it is rapidly becoming a recognizable hackathon template. Consumer shopping agents and payment controls are also adjacent to PayPal's own public agent strategy.

The thinner seam in the sampled work is **post-action correctness**: joining request intent, provider response, eventual webhook, retry/idempotency behavior, and merchant-visible state. This is a sample-based observation, not a quantified census. A product that makes this visible could serve developers/operators whose checkout/invoice/subscription systems call PayPal and need to resolve payment exceptions.

## White-space hypothesis

**PayPilot's candidate niche:** “When an AI-assisted or normal checkout payment has an uncertain outcome, prove what happened across PayPal and merchant-side records, then guide a safe next step.” This combines the product's existing code with the public strategic context of agentic commerce, without claiming an official PayPal partnership or inside knowledge.

For technical credibility, show a timeline with (a) internal order ID, (b) attempt/request idempotency key, (c) PayPal order/capture IDs, (d) sanitized webhook event IDs and received times, (e) authoritative refresh status, (f) confidence and evidence references, and (g) recovery policy. If these cannot be loaded from the real running backend, label the scenario synthetic and do not represent seeded events as PayPal transactions.
