# PayPilot demo walkthrough

## Submission video plan (under three minutes)

This is a recording outline, not evidence that an unverified segment has been completed. Keep the public simulation and the actual Sandbox integration visually and verbally separate. The hosted site is synthetic-only. Never imply its buttons call PayPal or an AI model.

### 0:00–0:20 — Problem and product

Open the [hosted judge demo](https://paypilot-hackathon-demo.onrender.com). Introduce PayPilot as an evidence-first payment operations workbench for resolving uncertain payment outcomes. Say that a timeout does not prove a payment failed, so retrying immediately can risk a duplicate capture.

### 0:20–1:20 — Synthetic walkthrough

1. Point out the **Guided simulation** banner and say: “This hosted walkthrough is synthetic. It does not call PayPal or an AI model, and it moves no money.”
2. Introduce the simulated agent mandate: “Find a warm desk lamp under $90 and buy one.” Explain that the illustrated $84 order is sample data.
3. Click **Run payment scenario**. Narrate the sequence as a simulation: intent is scoped, checkout is approved, the capture response is lost, and a completion event arrives later.
4. Pause at **Retry is paused until provider status is verified**. Explain that PayPilot preserves the unknown result rather than guessing that capture failed.
5. If you click **Verify with PayPal** in Demo mode, explicitly call it a simulated state transition. Do not describe it as a provider refresh. Show that the evidence/timeline reconciles the scenario without another capture.
6. Optional: run **Investigate** and state that the hosted walkthrough uses its synthetic, deterministic evidence only; it is not a live OpenAI response.

### 1:20–2:10 — Actual PayPal Sandbox evidence

Only include this segment from a separate local Live API recording after rechecking the environment. Show the actual local PayPilot order/payment state and the processed `PAYMENT.CAPTURE.COMPLETED` webhook for the verified USD $1 Sandbox transaction. Say: “Separately, we verified a buyer-approved PayPal Sandbox checkout, server-side capture, persisted paid status, and the matching verified webhook.” Do not show account emails, credentials, webhook secrets, full buyer details, or anything that could be mistaken for a live-money payment. If you cannot safely capture this screen, omit the visual and retain the exact verbal distinction.

### 2:10–2:40 — AI and limits (replace only after verification)

The OpenAI Responses API adapter is implemented, but a live model-backed investigation has not yet been verified. Do not say the product has a demonstrated live AI feature until issue #21 is complete. After a successful, redacted verification, show one investigation grounded in persisted order evidence, identify the model accurately, and explain that its tools are read-only while deterministic checks remain authoritative. If that evidence is not available, state plainly that AI runtime validation remains incomplete; the owner must decide whether the entry can proceed under the official rules.

### 2:40–2:55 — Close

End with: “PayPilot makes uncertainty visible, ties conclusions to payment evidence, and pauses risky recovery until provider state is known.” Identify the hosted URL as a synthetic judge demo and direct judges to the repository setup instructions for the database-backed Sandbox integration.

## Claims to keep precise

- Demo mode is a browser-local synthetic scenario. It calls neither PayPal nor an LLM and moves no money.
- Hosted Render mode is synthetic-only, has no database or payment/investigation API routes, and may take about a minute to wake after idle.
- A separate USD $1 PayPal Sandbox buyer checkout, server-side capture, persisted paid state, and verified processed webhook were completed on 9 November 2026. This is Sandbox evidence, not a live-money transaction.
- A live OpenAI investigation has not yet been verified. Do not present deterministic fallback output as model output.
- The prototype has no merchant authentication or tenant isolation. Never use production credentials or real merchant data.

## Reset to the known synthetic starting state

1. Select **Demo** mode, then click **Reset scenario** on the scenario card.
2. Confirm the status returns to **READY**, the action returns to **Run payment scenario**, and the selected sample is `PP-2026-0142`.
3. Run the walkthrough again from step 2 of the synthetic walkthrough above.

Reset is idempotent: clicking it repeatedly restores the same in-memory `demoBase` records. It also cancels an animation already in progress. The UI confirms **“Scenario reset. No provider calls were made.”** This action only replaces browser-local synthetic demo state; it does not call the Go API or modify PostgreSQL. Reloading the page also starts a fresh browser-local demo. Do not use this control as a database reset for any hosted environment.

**Verification (5 November 2026):** browser review clicked reset repeatedly and during a running scenario. Both paths returned to the same ready state, enabled the run action, and retained the selected seeded payment. The interrupted animation did not advance after reset. Issue #34 tracks this local synthetic reset.
