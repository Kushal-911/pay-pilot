# PayPilot

PayPilot is an evidence-grounded payment investigation and recovery platform for merchants using PayPal. The initial milestone implements a Go/Gin API, PostgreSQL order persistence, and the PayPal Sandbox checkout/capture server path.

## Run locally

1. Copy `.env.example` to `.env` and set PayPal Sandbox client credentials and checkout return/cancel URLs.
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

To apply new migrations to an existing local database, apply the numbered SQL migration once using `psql "$DATABASE_URL" -f migrations/002_reliability.sql`. Resetting the local Compose volume is another option, but deletes its data.

## API and project plan

- OpenAPI: `openapi.yaml`
- Postman: `postman/PayPilot.postman_collection.json`
- Milestones: `ROADMAP.md`
- Current status: `PROJECT_STATUS.md`
- Decisions: `DECISIONS.md`
- Session handoff: `HANDOFF.md`

The PayPal Sandbox API is not exercised by local tests; configure Sandbox credentials and accounts to complete the provider flow. Never use production credentials in local development.
