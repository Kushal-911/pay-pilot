# PayPilot

PayPilot is an evidence-grounded payment investigation and recovery platform for merchants using PayPal. The initial milestone implements a Go/Gin API, PostgreSQL order persistence, and the PayPal Sandbox checkout/capture server path.

## Merchant workbench

The same Go service serves an embedded, responsive merchant workbench at `http://localhost:8080/`. It includes an interactive, browser-only synthetic payment trace that demonstrates an agent-authorized purchase, an ambiguous capture timeout, delayed completion evidence, and a provider-state check before local reconciliation. The demo is explicitly labeled synthetic and does not call PayPal or move money. **Run the synthetic demo** is the reliable path for a recording without external accounts.

Choose **Live API** to read orders and investigations from the Go API and use PayPal Sandbox checkout, capture, refresh, and the configured investigation model. Live API requires the database and server; Sandbox actions additionally require PayPal Sandbox credentials. No automatic fallback disguises API failures as live data.

## Public hosted demo

The public deployment is intentionally synthetic-only. Set `PUBLIC_DEMO_MODE=true` to serve the browser workbench without a database or provider credentials; live API routes are not registered in this mode. The Render configuration is in [`render.yaml`](render.yaml), and setup, limitations, and verification steps are in [`DEPLOYMENT.md`](DEPLOYMENT.md). Do not add merchant data, PayPal credentials, or AI keys to the public demo service.

## Run locally

1. Copy `.env.example` to `.env` if you want to configure PayPal Sandbox client credentials, checkout return/cancel URLs, and `PAYPAL_WEBHOOK_ID`. To enable model-assisted investigation, also set `OPENAI_API_KEY` and `OPENAI_MODEL`; otherwise reports use deterministic record checks. These credentials are not needed to start the local app.
2. Start the database and API with `docker compose up --build -d`.
3. Open `http://localhost:8080/` for the workbench. `GET /livez` checks process liveness, while `GET /readyz` checks PostgreSQL readiness; `/health` remains a readiness alias.
4. To create a real Sandbox checkout, switch the workbench to **Live API** and create an order, or use `POST /api/v1/orders` followed by `POST /api/v1/orders/{id}/paypal-order`.
5. Open the returned `approval_url` in a browser and approve with a PayPal Sandbox buyer. Capture with the workbench or `POST /api/v1/payments/{payment_id}/capture`.

### Configure PayPal Sandbox

1. Sign in to the [PayPal Developer Dashboard](https://developer.paypal.com/dashboard/), switch to Sandbox, and open **Apps & Credentials**. Create a Sandbox app using a business Sandbox account, then copy its Client ID and secret into `.env` as `PAYPAL_CLIENT_ID` and `PAYPAL_CLIENT_SECRET`.
2. Under **Testing Tools → Sandbox Accounts**, use a personal account as the buyer and the business account as the seller. Open the checkout `approval_url` in a browser and sign in with the personal Sandbox credentials. PayPal's [Sandbox account guide](https://developer.paypal.com/sandbox-testing/accounts) describes how to find or create these test accounts.
3. For webhook verification, expose the local API at a reachable HTTPS URL, add `<https-url>/api/v1/webhooks/paypal` to the Sandbox app's webhooks, and subscribe to `PAYMENT.CAPTURE.COMPLETED`, `PAYMENT.CAPTURE.PENDING`, and `PAYMENT.CAPTURE.DENIED`. Save the registered webhook ID as `PAYPAL_WEBHOOK_ID` and restart the API. See PayPal's [webhook integration guide](https://developer.paypal.com/api/rest/webhooks/rest/).

Use Sandbox credentials and test accounts only. Keep `.env` out of version control, and never paste secrets into issues, logs, or the Postman collection.

Example order:

```json
{"merchant_reference":"demo-001","amount_minor":2500,"currency":"USD","description":"Demo order"}
```

If host ports 5432 or 8080 are occupied, set `POSTGRES_HOST_PORT` and/or `API_HOST_PORT` when starting Compose, for example `POSTGRES_HOST_PORT=55433 API_HOST_PORT=18082 docker compose up --build -d`.

PostgreSQL integration tests create a temporary schema, apply migrations 001–003, and drop the schema on completion. Run them explicitly against a disposable local database with `TEST_DATABASE_URL='postgres://paypilot:password@127.0.0.1:5432/paypilot?sslmode=disable' go test ./...`; without this variable, DB-gated tests are skipped.

`amount_minor` is an integer in the currency's minor unit. This initial implementation formats PayPal amounts with two decimal places, so use two-decimal currencies (for example USD) until currency exponent handling is expanded.

PostgreSQL initialization runs the numbered SQL files in `migrations/` on first creation of the Compose volume. Every mutation endpoint requires an `Idempotency-Key` header (8–128 characters); repeating the same key and request returns the stored response, while reusing it with a different request returns `409`. A processing record older than two minutes can be reclaimed after a crashed request; completed responses have no automatic expiry and must be retained for safe replay. Capture outcomes that cannot be confirmed are stored as `unknown`; call the refresh endpoint to retrieve the authoritative PayPal state before retrying. `GET /api/v1/orders/{id}/events` returns the selected order’s sanitized PayPal event history in received-time order, including processing status and correlation fields without raw webhook payloads.

To investigate an order, call `POST /api/v1/investigations` with an `order_id` and question, then inspect the persisted report and record-linked evidence. The report always includes deterministic cross-record checks. If OpenAI configuration is present, the Responses API can call four read-only evidence functions; its output is shown only after its evidence IDs and report fields validate. Provider errors or invalid output fall back to a labeled deterministic report. The investigation does not execute payment actions.

Configure the PayPal app webhook URL as `https://<your-host>/api/v1/webhooks/paypal`, subscribe it to `PAYMENT.CAPTURE.COMPLETED`, `PAYMENT.CAPTURE.PENDING`, and `PAYMENT.CAPTURE.DENIED`, and set its registered webhook ID as `PAYPAL_WEBHOOK_ID`. PayPal's webhook signature verification is called before an event is accepted; the event ID is unique and duplicate deliveries are acknowledged without reprocessing. See the [PayPal webhook integration guide](https://developer.paypal.com/api/rest/webhooks/rest/) for HTTPS and app registration requirements.

The current API has no merchant authentication or tenant isolation; use it only in a trusted local environment until those controls are implemented.

To apply new migrations to an existing local database, apply numbered SQL migrations once in order using `psql "$DATABASE_URL" -f migrations/002_reliability.sql`, `psql "$DATABASE_URL" -f migrations/003_investigations.sql`, and later numbered migration files. Resetting the local Compose volume is another option, but deletes its data.

## API and project plan

- OpenAPI: `openapi.yaml`
- Postman: `postman/PayPilot.postman_collection.json`

Import the Postman collection and set `baseUrl` to your local API. Send **Health**, **Create order**, then **List orders**, **List payment attempts**, and **List provider events**. The collection creates a run ID and saves the created `orderId`; clear its `runId` variable to create a fresh synthetic order. The PayPal order, capture, and webhook requests require Sandbox configuration and a buyer-approved order. Never capture an order before buyer approval; after an uncertain capture, call **Refresh payment** before considering another capture. The collection contains no credentials.
- Milestones: `ROADMAP.md`
- Current status: `PROJECT_STATUS.md`
- Decisions: `DECISIONS.md`
- Session handoff: `HANDOFF.md`
- Demo recording walkthrough: `DEMO_SCRIPT.md`

The PayPal Sandbox API is not exercised by local tests; configure Sandbox credentials and accounts to complete the provider flow. Never use production credentials in local development.
