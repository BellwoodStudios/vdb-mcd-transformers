CREATE TABLE maker.bite (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk       TEXT,
  urn       TEXT,
  ink       NUMERIC,
  art       NUMERIC,
  iArt      NUMERIC,
  tab       NUMERIC,
  nflip     NUMERIC,
  tx_idx    INTEGER NOT NUll,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx)
)