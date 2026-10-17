CREATE TABLE IF NOT EXISTS monitor_events (
    id TEXT PRIMARY KEY,
    chain TEXT NOT NULL,
    network TEXT NOT NULL,
    type TEXT NOT NULL,
    asset TEXT,
    wallet TEXT,
    amount DOUBLE PRECISION,
    price_usd DOUBLE PRECISION,
    liquidity_usd DOUBLE PRECISION,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    observed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_monitor_events_asset_time
    ON monitor_events (asset, observed_at DESC);

CREATE INDEX IF NOT EXISTS idx_monitor_events_type_time
    ON monitor_events (type, observed_at DESC);
