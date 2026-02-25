-- Migration 009: Seed services and link existing permissions

-- Auth Service (admin platform)
INSERT INTO services (id, tenant_id, name, slug, description, redirect_urls) VALUES
  ('e0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
   'Auth Service', 'auth-service', 'Plataforma de autenticacao e autorizacao centralizada',
   ARRAY['http://localhost:3003'])
ON CONFLICT DO NOTHING;

-- My Library
INSERT INTO services (id, tenant_id, name, slug, description, redirect_urls) VALUES
  ('e0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
   'My Library', 'my-library', 'Sistema de gerenciamento de biblioteca pessoal',
   ARRAY['http://localhost:3000'])
ON CONFLICT DO NOTHING;

-- Focus Hub
INSERT INTO services (id, tenant_id, name, slug, description, redirect_urls) VALUES
  ('e0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
   'Focus Hub', 'focus-hub', 'Sistema de produtividade com AI',
   ARRAY['http://localhost:3002'])
ON CONFLICT DO NOTHING;

-- Link existing permissions to services
UPDATE permissions SET service_id = 'e0000000-0000-0000-0000-000000000002'
  WHERE resource IN ('books', 'loans') AND service_id IS NULL;

UPDATE permissions SET service_id = 'e0000000-0000-0000-0000-000000000003'
  WHERE resource = 'tasks' AND service_id IS NULL;

UPDATE permissions SET service_id = 'e0000000-0000-0000-0000-000000000001'
  WHERE resource IN ('users', 'tenant') AND service_id IS NULL;
