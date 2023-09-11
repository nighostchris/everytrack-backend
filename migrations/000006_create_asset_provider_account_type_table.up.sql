CREATE TABLE IF NOT EXISTS everytrack_backend.asset_provider_account_type (
  id                UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  asset_provider_id UUID          REFERENCES everytrack_backend.asset_provider (id),
  name              TEXT          NOT NULL,
  created_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE(asset_provider_id, name)
);

CREATE TRIGGER asset_provider_account_type_updated_at
BEFORE UPDATE ON everytrack_backend.asset_provider_account_type
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
