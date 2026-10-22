# AI boundary and evidence design

## Where AI adds value

Operational records are structured but fragmented; the operator's question is free-form (“Did I charge them twice?” “What changed this morning?”). A language model can map that question into a small read-only investigation, summarize discrepancies in plain language, distinguish possible causes, and suggest which record to inspect next. The value is reduced search/synthesis effort and more accessible explanation.

The model is not the source of payment truth. It does not improve authoritative transaction state, sign webhook messages, decide accounting policy, or independently make payment actions safe.

## Correct trust boundary

```mermaid
sequenceDiagram
    participant U as Operator
    participant S as PayPilot investigation service
    participant DB as Payment records
    participant M as LLM
    participant PP as PayPal API
    U->>S: Natural-language question + local order ID
    S->>DB: Load scoped snapshot
    S->>M: Question + constrained read-only tools
    M->>S: Tool request: list attempts / compare records
    S->>DB: Validate target/order scope and execute query
    DB-->>S: Actual records with stable evidence IDs
    S-->>M: Structured tool result
    M-->>S: Draft findings + evidence references + hypotheses
    S->>S: Validate references and constrain schema/policy
    S-->>U: Facts, cited evidence, uncertainty, safe recommendation
    Note over S,PP: Separate refresh endpoint calls PayPal; LLM cannot call it directly
```

For this prototype the tool set is intentionally small: `get_payment`, `get_attempt_history`, `list_provider_events`, and `compare_order_records`. Tools accept an order ID and enforce scope; no arbitrary SQL or user-provided URLs. Validate arguments, bound result size, time out model calls, and fall back cleanly. Do not send unnecessary personal/payment data to the model.

## Evidence representation

Each fact should include a stable evidence reference, source type, record ID, observation time, and a short relevant value. Example:

```json
{
  "claim": "PayPal reports the capture completed.",
  "classification": "confirmed",
  "evidence": [
    {"ref":"provider_refresh:ord_abc:2026-10-19T10:20:00Z","field":"status","value":"COMPLETED"}
  ]
}
```

In the actual app, use real IDs from persisted data; the example is illustrative. At display time, the backend must confirm every reference belongs to the selected order and exists. Never accept a model-supplied event ID or amount as evidence without resolving it against storage/provider output.

### Three kinds of language

- **Confirmed fact:** directly present in verified provider response or persisted record and cited.
- **Hypothesis:** a plausible interpretation, clearly labeled and not presented as a status.
- **Unknown:** evidence missing or conflicting. Say what additional read/record would resolve it.

Confidence should mean confidence in the diagnostic explanation, with a reason (e.g. “high because provider refresh and completed webhook agree on capture ID and amount”), not an unexplained scalar that sounds statistically calibrated.

## AI safety and error behavior

- Model output cannot update order/attempt status or send payment/refund/dispute mutations.
- The backend owns auth, policy and state transitions.
- Read-only tools are order-scoped and fixed; no direct DB access for the model.
- Tool outputs are untrusted data; merchant descriptions/events may contain hostile prompt text. Treat as quoted data, not instructions.
- Do not expose secrets, full raw webhook bodies, customer payment data, or credentials.
- Validate structured output, evidence IDs, enum values, lengths and claim classifications.
- On timeout/invalid output/provider outage, return deterministic findings labeled as AI unavailable/fallback; never mark payment operations failed because the explanation service failed.
- Keep model version, prompt/tool schema versions, investigation time, and evidence snapshot ID for reproducibility if persistence/privacy policy permits.

## Why not use a bespoke ML model?

The MVP's task is constrained retrieval and explanation over a small set of records, not prediction with labeled outcomes. A large custom model would require reliable labels, representative data, evaluation and ongoing governance. Deterministic invariant checks should carry payment correctness. A hosted LLM is suitable for human-readable synthesis if it is bounded and its contribution is tested against a deterministic baseline.

## Demonstrate AI value honestly

Show the same real evidence two ways: raw IDs/records and a natural-language investigation. The LLM should explain a non-obvious relationship (e.g. an API timeout, later completion event and local pending state) and cite all three. Include a model-unavailable path and a deliberately invalid evidence response in tests. Do not use AI to manufacture fraud scores or payment outcomes.
