CREATE TABLE IF NOT EXISTS everytrack_backend.exchange_rate (
  id                  UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  base_currency_id    UUID          REFERENCES everytrack_backend.currency (id),
  target_currency_id  UUID          REFERENCES everytrack_backend.currency (id),
  rate                TEXT          NOT NULL,
  created_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE(base_currency_id, target_currency_id)
);

CREATE TRIGGER exchange_rate_updated_at
BEFORE UPDATE ON everytrack_backend.exchange_rate
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
