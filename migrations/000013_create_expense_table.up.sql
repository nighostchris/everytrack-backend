CREATE TABLE IF NOT EXISTS everytrack_backend.expense (
  id            UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  client_id     UUID          REFERENCES everytrack_backend.client (id),
  account_id    UUID          REFERENCES everytrack_backend.account (id),
  currency_id   UUID          REFERENCES everytrack_backend.currency (id),
  category      TEXT          NOT NULL,
  amount        TEXT          NOT NULL,
  remarks       TEXT,
  executed_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TRIGGER expense_updated_at
BEFORE UPDATE ON everytrack_backend.expense
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
