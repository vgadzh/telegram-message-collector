CREATE TABLE users
(
  id            BIGSERIAL PRIMARY KEY,
  login         TEXT        NOT NULL UNIQUE,
  password_hash TEXT        NOT NULL,
  role          TEXT        NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);