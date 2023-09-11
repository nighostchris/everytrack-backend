CREATE TABLE IF NOT EXISTS everytrack_backend.asset_provider (
  id          UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  name        TEXT          UNIQUE NOT NULL,
  icon        TEXT          NOT NULL,
  type        VARCHAR(10)   NOT NULL,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TRIGGER asset_provider_updated_at
BEFORE UPDATE ON everytrack_backend.asset_provider
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
