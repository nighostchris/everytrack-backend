CREATE TABLE IF NOT EXISTS everytrack_backend.account (
  id                              UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  client_id                       UUID          REFERENCES everytrack_backend.client (id),
  asset_provider_account_type_id  UUID          REFERENCES everytrack_backend.asset_provider_account_type (id),
  currency_id                     UUID          REFERENCES everytrack_backend.currency (id),
  balance                         TEXT          NOT NULL,
  created_at                      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at                      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE(client_id, asset_provider_account_type_id)
);

CREATE TRIGGER account_updated_at
BEFORE UPDATE ON everytrack_backend.account
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
