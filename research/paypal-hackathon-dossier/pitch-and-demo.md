# Pitch, demo, and submission outline

## Product name and one-liner

**PayPilot — Transaction Truth for Agentic Commerce**  
When a PayPal payment outcome is unclear, PayPilot reconstructs the evidence and guides a safe next step—so a merchant doesn't retry a payment just because their own request timed out.

Keep “agentic commerce” as context, not a claim that the current prototype is PayPal's official agent platform. If the finished UI has no actual agent-originated action, use “AI-enabled commerce” in the tagline.

## 30-second verbal pitch

“An API timeout doesn't tell a merchant whether the payment failed. It may have succeeded while the merchant's app missed the response. Retrying blindly risks charging the customer twice; waiting blindly strands an order. PayPilot correlates the merchant order, PayPal capture attempt, verified PayPal webhooks, and a fresh provider status. It shows a timestamped evidence trail, explains mismatches with linked records, and blocks a retry until backend rules say it's safe. AI helps investigate; PayPal's verified response remains the authority.”

## Three-minute demo script (target 2:40)

| Time | Screen/action | Narration and proof |
|---|---|---|
| 0:00–0:20 | Show exception queue with one pending/unknown payment, a clear “Sandbox demo” badge | “This is the failure mode: PayPal may have completed the capture, but the merchant's request timed out.” State clearly if using seeded data. |
| 0:20–0:45 | Open order; show merchant status, amount/currency, attempt ID, PayPal order/capture IDs | “The merchant record alone is not authoritative. Here are the correlated identifiers and event history.” |
| 0:45–1:10 | Click “Refresh PayPal status”; show an actual Sandbox API response and status transition | “We query PayPal before deciding what to do; a timeout is not interpreted as failure.” If no real call is working, do not fake this segment. |
| 1:10–1:45 | Open timeline; show provider capture and verified webhook with duplicate/out-of-order badge or a clearly labeled synthetic injection | “Duplicate delivery is idempotent; a late pending event doesn't downgrade a confirmed completion.” |
| 1:45–2:20 | Run investigation; show concise finding, evidence refs that navigate to rows, confidence explanation, hypotheses separated from confirmed facts | “The model can retrieve only scoped read-only evidence. The backend validates the references; deterministic policy controls status.” |
| 2:20–2:40 | Show safe action/reason and end state | “PayPilot recommends no second capture because PayPal confirms completion. That is the business outcome: resolve the exception without asking the customer to pay twice.” |

Avoid spending time on architecture diagrams, prompt engineering, or a generic chatbot conversation. The recording must show the project actually running on the device for which it was built and stay under three minutes.

## Devpost text skeleton

### Inspiration
Payment APIs are asynchronous distributed systems. When a response is lost, the merchant's internal state can diverge from PayPal's state. Operators then search logs and dashboards and may risk an unsafe retry. PayPilot investigates this narrow exception path.

### What it does
Describe only implemented functions: merchant order/capture, PayPal Sandbox calls, webhook verification/deduplication, deterministic state rules, investigation evidence and UI. Mark synthetic scenarios.

### How we built it
Go/Gin service, PostgreSQL, PayPal Orders v2 and webhooks, [frontend stack once chosen], OpenAI Responses API with four read-only tools. AI has no database credentials, arbitrary SQL, or payment-mutation tool. Mention deterministic fallback and backend evidence validation.

### What makes it different
Not a generic agent that can call a payment API and not an LLM that declares payments successful. It focuses on the uncertain boundary between a request, the provider's authoritative outcome, webhook delivery, and merchant persistence. Show the behavior rather than asserting uniqueness.

### Challenges / what's next
Be candid about sandbox verification, current environment constraints, no auth/tenant isolation, supported currency scope, and what remains a prototype. Include tests actually run.

## Visual design principles

- Use status names with precise meaning: “Unknown locally,” “PayPal reports COMPLETED,” “Webhook received/verified,” rather than a single ambiguous green Paid badge.
- Every factual sentence in an AI finding should expose an evidence reference.
- Mark synthetic data persistently, not only in README.
- Make amount/currency, provider IDs, event times, and action availability visible without scrolling into raw JSON.
- A failed AI call should leave the order and deterministic checks usable.

## Pitch claim audit

Before recording, require evidence for each claim. “Prevents duplicate capture” means tests and state policy block the second capture in the ambiguous state. “Grounded” means invalid evidence IDs are rejected in code and test. “Real PayPal integration” means an actual Sandbox flow was executed and recorded. “Saves time” is a value hypothesis until a measured user test. “PayPal wants this” should not be claimed; say the project is aligned with public developer and trust direction.
