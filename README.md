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

PostgreSQL initialization runs `migrations/001_initial.sql` on the first creation of the Compose volume. To apply changed initial migrations locally, reset the local database volume deliberately with `docker compose down -v` and restart.

## API and project plan

- OpenAPI: `openapi.yaml`
- Postman: `postman/PayPilot.postman_collection.json`
- Milestones: `ROADMAP.md`
- Current status: `PROJECT_STATUS.md`
- Decisions: `DECISIONS.md`
- Session handoff: `HANDOFF.md`

The PayPal Sandbox API is not exercised by local tests; configure Sandbox credentials and accounts to complete the provider flow. Never use production credentials in local development.
