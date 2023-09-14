CREATE TABLE IF NOT EXISTS everytrack_backend.account_stock (
  id            UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  account_id    UUID          REFERENCES everytrack_backend.account (id),
  stock_id      UUID          REFERENCES everytrack_backend.stock (id),
  unit          TEXT          NOT NULL,
  cost          TEXT          NOT NULL,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE(account_id, stock_id)
);

CREATE TRIGGER account_stock_updated_at
BEFORE UPDATE ON everytrack_backend.account_stock
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
