CREATE TABLE provider_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_event_id TEXT NOT NULL UNIQUE,
    order_id UUID REFERENCES orders(id),
    event_type TEXT NOT NULL,
    paypal_order_id TEXT,
    paypal_capture_id TEXT,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    processing_status TEXT NOT NULL DEFAULT 'received'
        CHECK (processing_status IN ('received','processed','failed')),
    received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);
CREATE INDEX provider_events_order_idx ON provider_events(paypal_order_id, received_at DESC);
CREATE INDEX provider_events_capture_idx ON provider_events(paypal_capture_id, received_at DESC);
CREATE INDEX provider_events_local_order_idx ON provider_events(order_id, received_at DESC);

CREATE TABLE investigations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id),
    question TEXT NOT NULL,
    finding TEXT NOT NULL,
    confidence DOUBLE PRECISION NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    recommendation TEXT NOT NULL CHECK (recommendation IN ('no_action','refresh_provider_status','review_records','manual_review')),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open','resolved','dismissed')),
    evidence JSONB NOT NULL DEFAULT '[]'::jsonb,
    report JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX investigations_order_created_idx ON investigations(order_id, created_at DESC);
