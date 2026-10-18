# PayPilot

PayPilot is an evidence-grounded payment investigation and recovery platform for merchants using PayPal. The initial milestone implements a Go/Gin API, PostgreSQL order persistence, and the PayPal Sandbox checkout/capture server path.

## Run locally

1. Copy `.env.example` to `.env` and set PayPal Sandbox client credentials, checkout return/cancel URLs, and `PAYPAL_WEBHOOK_ID`. To enable model-assisted investigation, also set `OPENAI_API_KEY` and `OPENAI_MODEL`; otherwise reports use deterministic record checks.
2. Start the database and API with `docker compose up --build`.
3. Check `GET http://localhost:8080/health`.
4. Create an order with `POST /api/v1/orders`, then create a PayPal checkout with `POST /api/v1/orders/{id}/paypal-order`.
5. Open the returned `approval_url` in a browser and approve with a PayPal Sandbox buyer. Capture with `POST /api/v1/payments/{payment_id}/capture`.

Example order:

```json
{"merchant_reference":"demo-001","amount_minor":2500,"currency":"USD","description":"Demo order"}
```

`amount_minor` is an integer in the currency's minor unit. This initial implementation formats PayPal amounts with two decimal places, so use two-decimal currencies (for example USD) until currency exponent handling is expanded.

PostgreSQL initialization runs the numbered SQL files in `migrations/` on first creation of the Compose volume. Every mutation endpoint requires an `Idempotency-Key` header (8–128 characters); repeating the same key and request returns the stored response, while reusing it with a different request returns `409`. Capture outcomes that cannot be confirmed are stored as `unknown`; call the refresh endpoint to retrieve the authoritative PayPal state before retrying.

To investigate an order, call `POST /api/v1/investigations` with an `order_id` and question, then inspect the persisted report and record-linked evidence. The report always includes deterministic cross-record checks. If OpenAI configuration is present, the Responses API can call four read-only evidence functions; its output is shown only after its evidence IDs and report fields validate. Provider errors or invalid output fall back to a labeled deterministic report. The investigation does not execute payment actions.

Configure the PayPal app webhook URL as `https://<your-host>/api/v1/webhooks/paypal`, subscribe it to `PAYMENT.CAPTURE.COMPLETED`, `PAYMENT.CAPTURE.PENDING`, and `PAYMENT.CAPTURE.DENIED`, and set its registered webhook ID as `PAYPAL_WEBHOOK_ID`. PayPal's webhook signature verification is called before an event is accepted; the event ID is unique and duplicate deliveries are acknowledged without reprocessing. See the [PayPal webhook integration guide](https://developer.paypal.com/api/rest/webhooks/rest/) for HTTPS and app registration requirements.

The current API has no merchant authentication or tenant isolation; use it only in a trusted local environment until those controls are implemented.

To apply new migrations to an existing local database, apply numbered SQL migrations once in order using `psql "$DATABASE_URL" -f migrations/002_reliability.sql`, `psql "$DATABASE_URL" -f migrations/003_investigations.sql`, and later numbered migration files. Resetting the local Compose volume is another option, but deletes its data.

## API and project plan

- OpenAPI: `openapi.yaml`
- Postman: `postman/PayPilot.postman_collection.json`
- Milestones: `ROADMAP.md`
- Current status: `PROJECT_STATUS.md`
- Decisions: `DECISIONS.md`
- Session handoff: `HANDOFF.md`

The PayPal Sandbox API is not exercised by local tests; configure Sandbox credentials and accounts to complete the provider flow. Never use production credentials in local development.
