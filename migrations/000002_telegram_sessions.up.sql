CREATE TABLE sessions
(
  id               BIGSERIAL PRIMARY KEY,
  user_id          BIGINT      NOT NULL,
  telegram_user_id BIGINT,
  phone            TEXT        NOT NULL,
  phone_hash       TEXT,
  enc_session      TEXT        NOT NULL,
  status           TEXT        NOT NULL,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at       TIMESTAMPTZ,

  CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id),
  CONSTRAINT session_status_check
    CHECK (
      status IN (
                 'SENDING_CODE',
                 'PROCESSING_QR',
                 'CODE_REQUIRED',
                 'VERIFYING_CODE',
                 'PASSWORD_REQUIRED',
                 'VERIFYING_PASSWORD',
                 'ACTIVE',
                 'INACTIVE',
                 'FAILED'
        )
      )
);
CREATE UNIQUE INDEX uq_session_tg_user ON sessions (telegram_user_id) WHERE deleted_at IS NULL AND telegram_user_id > 0;
CREATE UNIQUE INDEX uq_session_phone_user_status ON sessions (user_id, phone) WHERE deleted_at IS NULL AND status IN ('SENDING_CODE', 'PROCESSING_QR');
