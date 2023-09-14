CREATE TABLE IF NOT EXISTS everytrack_backend.country (
  id          UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  name        TEXT          NOT NULL,
  code        VARCHAR(2)    UNIQUE NOT NULL,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TRIGGER country_updated_at
BEFORE UPDATE ON everytrack_backend.country
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
