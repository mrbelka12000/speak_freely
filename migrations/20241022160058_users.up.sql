BEGIN;

ALTER TABLE users ADD COLUMN IF NOT EXISTS external_id TEXT default '';

COMMIT;