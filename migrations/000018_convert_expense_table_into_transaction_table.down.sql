ALTER TABLE everytrack_backend.transaction
DROP COLUMN IF EXISTS income;

ALTER TRIGGER transaction_updated_at
ON everytrack_backend.transaction
RENAME TO expense_updated_at;

ALTER INDEX everytrack_backend.transaction_pkey RENAME TO expense_pkey;

ALTER TABLE everytrack_backend.transaction
RENAME CONSTRAINT transaction_currency_id_fkey TO expense_currency_id_fkey;

ALTER TABLE everytrack_backend.transaction
RENAME CONSTRAINT transaction_client_id_fkey TO expense_client_id_fkey;

ALTER TABLE everytrack_backend.transaction
RENAME CONSTRAINT transaction_account_id_fkey TO expense_account_id_fkey;

ALTER TABLE everytrack_backend.transaction
RENAME TO expense;
