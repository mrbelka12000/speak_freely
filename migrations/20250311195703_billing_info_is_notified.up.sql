BEGIN;

ALTER TABLE billing_info ADD COLUMN IF NOT EXISTS is_notified boolean default true;

COMMIT;