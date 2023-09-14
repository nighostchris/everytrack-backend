CREATE TABLE IF NOT EXISTS everytrack_backend.stock (
  id              UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  country_id      UUID          REFERENCES everytrack_backend.country (id),
  currency_id     UUID          REFERENCES everytrack_backend.currency (id),
  name            TEXT          NOT NULL,
  ticker          VARCHAR(20)   NOT NULL,
  current_price   TEXT          NOT NULL,
  created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE(country_id, ticker)
);

CREATE TRIGGER stock_updated_at
BEFORE UPDATE ON everytrack_backend.stock
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
