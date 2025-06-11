CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS trade (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    close_time TEXT,
    trade_date TIMESTAMPTZ NOT NULL,
    instrument_code TEXT NOT NULL,
    trade_price NUMERIC(15, 2) NOT NULL,
    trade_quantity INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted BOOLEAN DEFAULT false
);
