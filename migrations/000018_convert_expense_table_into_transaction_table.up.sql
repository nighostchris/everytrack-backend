ALTER TABLE everytrack_backend.expense
RENAME TO "transaction";

ALTER INDEX everytrack_backend.expense_pkey RENAME TO transaction_pkey;

ALTER TABLE everytrack_backend.transaction
RENAME CONSTRAINT expense_account_id_fkey TO transaction_account_id_fkey;

ALTER TABLE everytrack_backend.transaction
RENAME CONSTRAINT expense_client_id_fkey TO transaction_client_id_fkey;

ALTER TABLE everytrack_backend.transaction
RENAME CONSTRAINT expense_currency_id_fkey TO transaction_currency_id_fkey;

ALTER TRIGGER expense_updated_at
ON everytrack_backend.transaction
RENAME TO transaction_updated_at;

ALTER TABLE everytrack_backend.transaction
ADD COLUMN income BOOLEAN NOT NULL DEFAULT FALSE;
