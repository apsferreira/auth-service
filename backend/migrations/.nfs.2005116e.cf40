-- Migration 007: Create services table for dynamic application registration

CREATE TABLE services (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  description TEXT,
  redirect_urls TEXT[],
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(tenant_id, slug)
);

CREATE INDEX idx_services_tenant_id ON services(tenant_id);
CREATE INDEX idx_services_slug ON services(slug);

COMMENT ON TABLE services IS 'Aplicacoes/servicos registrados na plataforma de autenticacao';
COMMENT ON COLUMN services.slug IS 'Identificador unico do servico (ex: my-library, focus-hub)';
COMMENT ON COLUMN services.redirect_urls IS 'URLs permitidas para redirect pos-login';
