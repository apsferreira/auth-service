-- Migration 008: Add service_id to permissions for service-scoped authorization

ALTER TABLE permissions ADD COLUMN service_id UUID REFERENCES services(id) ON DELETE CASCADE;

-- Drop old global unique constraint on name
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS permissions_name_key;

-- New unique: permission name unique per service (NULL service_id treated as global)
CREATE UNIQUE INDEX idx_permissions_service_name
  ON permissions(COALESCE(service_id, '00000000-0000-0000-0000-000000000000'), name);

CREATE INDEX idx_permissions_service_id ON permissions(service_id);

COMMENT ON COLUMN permissions.service_id IS 'Servico ao qual esta permissao pertence (NULL = permissao global)';
