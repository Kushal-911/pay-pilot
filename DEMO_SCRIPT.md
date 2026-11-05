# PayPilot demo walkthrough

## 90-second recording

1. Open the workbench at `http://localhost:8080/` (or the active local preview URL) with **Demo** selected. Point out the **Synthetic Scenario** label and the banner stating that no money moves.
2. Introduce the agent mandate: “Find a warm desk lamp under $90 and buy one.” Explain that the $90 ceiling is part of the purchase intent, and the illustrated order is $84.
3. Click **Run payment scenario**. As the four steps animate, narrate: “The agent intent is scoped; checkout is approved; capture times out, so the system records an unknown result instead of assuming failure; then the completion event arrives late.”
4. Pause on **Retry is paused until provider status is verified**. Explain that retrying after a lost response can create duplicate charges.
5. Click **Verify with PayPal**. Show the provider refresh, matching capture references, and the final timeline entry **Merchant order reconciled**. Call out that the order leaves the exception queue and that no second capture was sent.
6. Click **Investigate** if time allows. Explain that the guided analysis is local and deterministic, while Live API mode uses the existing Go investigation endpoint and evidence validation.
7. End with the evidence-first principle: “A timeout is not a payment outcome. PayPilot waits for provider truth, then repairs local state once.”

## Claims to keep precise

- Demo mode is a browser-local synthetic scenario. It calls neither PayPal nor an LLM and moves no money.
- Live API mode is wired to the Go API. A real Sandbox checkout requires a running PostgreSQL/API stack and PayPal Sandbox credentials.
- The current prototype has no merchant authentication or tenant isolation. Keep it local and use only Sandbox credentials.

## Reset to the known synthetic starting state

1. Select **Demo** mode, then click **Reset scenario** on the scenario card.
2. Confirm the status returns to **READY**, the action returns to **Run payment scenario**, and the selected sample is `PP-2026-0142`.
3. Run the walkthrough again from step 2 of the recording script.

Reset is idempotent: clicking it repeatedly restores the same in-memory `demoBase` records. It also cancels an animation already in progress. The UI confirms **“Scenario reset. No provider calls were made.”** This action only replaces browser-local synthetic demo state; it does not call the Go API or modify PostgreSQL. Reloading the page also starts a fresh browser-local demo. Do not use this control as a database reset for any hosted environment.

**Verification (5 November 2026):** browser review clicked reset repeatedly and clicked reset during a running scenario. Both paths returned to the same ready state, enabled the run action, and retained the selected seeded payment. The interrupted animation did not advance after reset. Issue #34 tracks this local synthetic reset; a hosted environment still needs separate isolation and reset review before deployment.
