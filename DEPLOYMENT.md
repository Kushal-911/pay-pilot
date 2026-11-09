# Public demo deployment

PayPilot's public demo is an isolated, browser-only simulation. It serves the synthetic workbench and health checks; it does not connect to PostgreSQL, PayPal, OpenAI, or any merchant data. The public mode does not register order, payment, webhook, or investigation API routes.

## Deploy on Render

The repository includes a Render Blueprint at [`render.yaml`](render.yaml). It creates a free Docker web service in Singapore with `PUBLIC_DEMO_MODE=true` and `/readyz` as its health check.

1. Sign in to Render and choose **New → Blueprint**.
2. Connect the public `Kushal-911/pay-pilot` GitHub repository and select `main`.
3. Review the Blueprint and create the service. It requires no database, API keys, or other secrets.
4. Wait for the deploy to become live, then open the service's `onrender.com` URL.

The Render Free plan is suitable for a hackathon preview, with platform limits. A free web service spins down after 15 minutes without inbound traffic; its first request after sleeping can take about a minute. Free services also have monthly instance-hour, build, and bandwidth limits. Review [Render's Free plan limits](https://render.com/docs/free) before relying on the URL during judging.

## Verify the public service

Replace `<host>` with the Render hostname:

```sh
curl -i https://<host>/
curl -i https://<host>/livez
curl -i https://<host>/readyz
curl -i https://<host>/api/v1/orders
```

The workbench and both health checks should return `200`. The orders route should return `404`. Payment, capture, webhook, and investigation API routes are also absent. The page identifies its records as synthetic and hides Live API and new-payment controls.

## Local public-mode check

Run the public mode without Docker or a database:

```sh
PUBLIC_DEMO_MODE=true PORT=8080 GIN_MODE=release go run ./cmd/api
```

Then open `http://localhost:8080/`. This mode has no provider or database configuration. For the full PayPal Sandbox flow, use the normal local Compose setup in the README with your own Sandbox credentials; do not put those credentials into Render.

## Rollback and scope

Use the Render service's **Events** or **Deploys** page to roll back to a previous deploy if a new version fails. The demo has no server-side payment data to reset or preserve; all records shown are synthetic and live only in the browser. Keep the local, database-backed Live API mode separate from this public service. Never use this deployment for production PayPal credentials or real customer data.
