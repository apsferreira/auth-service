-- Migration 001: Create tenants table

CREATE TABLE tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(100) NOT NULL UNIQUE,
  plan VARCHAR(50) NOT NULL DEFAULT 'free',
  settings JSONB DEFAULT '{}'::jsonb,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tenants_slug ON tenants(slug);

COMMENT ON TABLE tenants IS 'Tenants (organizations/clients)';
COMMENT ON COLUMN tenants.slug IS 'Unique slug for tenant (e.g., apsferreira)';
COMMENT ON COLUMN tenants.plan IS 'Tenant plan: free, premium, enterprise';

