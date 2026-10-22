# Business case and validation plan

## Business problem, stated narrowly

When a merchant's local payment state differs from PayPal's state, the work to diagnose it crosses application logs, database rows, provider dashboards, and webhook history. During an ambiguous capture outcome, the merchant must balance two costly risks: retrying and possibly charging twice, or waiting and delaying fulfillment. PayPilot proposes a single evidence trail and a safe reconciliation workflow.

This is a product hypothesis, not a validated incidence rate. No internal PayPal operational data was accessed and no merchant interviews have yet been reported.

## Value chain

| Product capability | Immediate operator value | Potential business effect to test |
|---|---|---|
| Correlated order/attempt/provider timeline | Fewer manual cross-system searches | Lower time-to-resolution and support handling cost |
| Provider refresh before retry | Evidence before consequential action | Fewer duplicate captures and payment confusion |
| Duplicate/out-of-order event handling | Stable local order state under retries | Fewer fulfillment errors and reconciliation tickets |
| Evidence-cited investigation summary | Understandable explanation without raw log expertise | More operators can self-serve; less escalation to engineering |
| AI fallback to deterministic findings | Basic operations remain available if model is down | Lower operational dependency on model availability |

For PayPal, possible value is a healthier developer ecosystem and fewer avoidable merchant integration failures; measure developer setup time, Sandbox completion, and support needs before making that claim strongly.

## Why pay for AI here?

Do not say “because AI is the future.” The user pain is interpretation across records and translating a question into relevant checks. AI is useful if it reduces time and expertise needed to find the linked facts. Payment status remains deterministic because the costs of a model hallucination are unacceptable. This division is economically sensible: use probabilistic language understanding where it saves human effort; use provider APIs and code where state correctness matters.

## What is not claimed

- Not a fraud detector, chargeback win-rate predictor, guaranteed chargeback reduction, or autonomous payment executor.
- Not a PayPal product, partner, approved integration, or endorsed architecture.
- Not proof of production financial return or average merchant pain.
- Not safe for production merchants until authentication, tenant isolation, privacy, security and reliability work is completed.

## Validation plan

### Discovery interviews (before adding scope)

Interview five small merchants, e-commerce developers, or support/finance operators. Ask:

1. Tell me about the last PayPal payment whose status looked wrong.
2. What did you check and in what order?
3. How long did it take, who else got involved, and what was the outcome?
4. Has a retry, delayed fulfillment, or duplicate request caused a customer issue?
5. What evidence would you need before trusting a suggested next step?
6. Which systems contain the useful records (store platform, application logs, PayPal dashboard, shipping, support)?

Record anonymized quotes only with consent. Do not ask participants to share real payment credentials or customer data.

### Prototype usability test

Give three people a timed scenario and let them use the UI without coaching. Measure: time to find provider status, whether they can tell fact from hypothesis, whether they understand why retry is blocked/allowed, whether evidence refs are meaningful, and wrong action rate. Fix confusing status language before visual polish.

### Pilot metrics

Baseline and compare a small number of exception cases:

- Time-to-resolution median and 90th percentile.
- Manual navigation/search steps per exception.
- Operator agreement with evidence-backed diagnosis.
- Evidence reference click-through and correction rate.
- Incorrect state transition or duplicate action (must be zero in prototype tests).
- Model fallback availability and investigation latency.

For real-world ROI, define labor cost, exception volume, customer support rate and duplicate loss before estimating savings. A prototype should report measured values, not projected numbers presented as achieved results.
