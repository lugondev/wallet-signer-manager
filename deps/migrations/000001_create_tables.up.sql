BEGIN;

CREATE TABLE IF NOT EXISTS wallets
(
    pk                    SERIAL PRIMARY KEY,
    store_id              TEXT                                           NOT NULL,
    public_key            BYTEA                                          NOT NULL,
    pubkey                text                                           NOT NULL,
    compressed_public_key BYTEA                                          NOT NULL,
    tags                  JSONB,
    auth                  JSONB,
    disabled              BOOLEAN     default false,
    created_at            TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    updated_at            TIMESTAMPTZ DEFAULT (now() at time zone 'utc') NOT NULL,
    deleted_at            TIMESTAMPTZ,
    UNIQUE (pubkey, store_id)
);

COMMIT;
