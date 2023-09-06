CREATE OR REPLACE FUNCTION on_update_timestamp() RETURNS trigger
LANGUAGE 'plpgsql'
AS $$
  BEGIN
    NEW.updated_at = now();
    RETURN NEW;
  END;
$$;
