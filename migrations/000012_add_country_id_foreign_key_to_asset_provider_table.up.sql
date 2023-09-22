ALTER TABLE everytrack_backend.asset_provider
ADD COLUMN country_id UUID REFERENCES everytrack_backend.country (id);

ALTER TABLE everytrack_backend.asset_provider
DROP CONSTRAINT IF EXISTS asset_provider_name_key;

ALTER TABLE everytrack_backend.asset_provider
ADD CONSTRAINT asset_provider_name_country_id_key UNIQUE(name, country_id);
