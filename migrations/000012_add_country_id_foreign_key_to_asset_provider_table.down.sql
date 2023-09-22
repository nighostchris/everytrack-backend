ALTER TABLE everytrack_backend.asset_provider
DROP COLUMN IF EXISTS country_id;

ALTER TABLE everytrack_backend.asset_provider
ADD CONSTRAINT asset_provider_name_key UNIQUE(name);
