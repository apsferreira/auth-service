DROP INDEX IF EXISTS idx_permissions_service_name;
DROP INDEX IF EXISTS idx_permissions_service_id;
ALTER TABLE permissions DROP COLUMN IF EXISTS service_id;
ALTER TABLE permissions ADD CONSTRAINT permissions_name_key UNIQUE (name);
