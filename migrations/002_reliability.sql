ALTER TABLE payment_attempts DROP CONSTRAINT payment_attempts_status_check;
ALTER TABLE payment_attempts ADD CONSTRAINT payment_attempts_status_check
    CHECK (status IN ('created','capturing','unknown','completed','failed'));

CREATE TABLE idempotency_keys (
    key TEXT PRIMARY KEY,
    request_hash TEXT NOT NULL,
    state TEXT NOT NULL CHECK (state IN ('processing','completed')),
    response_status INTEGER,
    stored_response JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK ((state = 'processing' AND response_status IS NULL AND stored_response IS NULL)
        OR (state = 'completed' AND response_status IS NOT NULL AND stored_response IS NOT NULL))
);
CREATE INDEX idempotency_keys_updated_idx ON idempotency_keys(updated_at);
