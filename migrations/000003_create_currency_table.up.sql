CREATE TABLE IF NOT EXISTS everytrack_backend.currency (
  id          UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  ticker      TEXT          UNIQUE NOT NULL,
  symbol      VARCHAR(5)    NOT NULL,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TRIGGER currency_updated_at
BEFORE UPDATE ON everytrack_backend.currency
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
