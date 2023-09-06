CREATE TABLE IF NOT EXISTS everytrack_backend.client (
  id          UUID          DEFAULT gen_random_uuid() PRIMARY KEY,
  email       TEXT          UNIQUE NOT NULL,
  username    VARCHAR(20)   NOT NULL,
  password    TEXT          NOT NULL,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE TRIGGER client_updated_at
BEFORE UPDATE ON everytrack_backend.client
FOR EACH ROW
EXECUTE PROCEDURE on_update_timestamp();
