CREATE TABLE IF NOT EXISTS everytrack_backend.future_payment (
  id              UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  client_id       UUID          REFERENCES everytrack_backend.client (id),
  account_id      UUID          REFERENCES everytrack_backend.account (id),
  currency_id     UUID          REFERENCES everytrack_backend.currency (id),
  name            TEXT          NOT NULL,
  amount          TEXT          NOT NULL,
  income          BOOLEAN       NOT NULL DEFAULT FALSE,
  rolling         BOOLEAN       NOT NULL DEFAULT FALSE,
  frequency       BIGINT,       
  remarks         TEXT,
  scheduled_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TRIGGER future_payment_updated_at
BEFORE UPDATE ON everytrack_backend.future_payment
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
