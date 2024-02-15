CREATE TABLE IF NOT EXISTS everytrack_backend.cash (
  id            UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  client_id     UUID          REFERENCES everytrack_backend.client (id),
  currency_id   UUID          REFERENCES everytrack_backend.currency (id),
  amount        TEXT          NOT NULL,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE(client_id, currency_id)
);

CREATE TRIGGER cash_updated_at
BEFORE UPDATE ON everytrack_backend.cash
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
