-- Migration 006: Seed data for development/testing

-- Default tenant
INSERT INTO tenants (id, name, slug, plan, settings)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'APS Ferreira',
    'apsferreira',
    'premium',
    '{"max_users": 100}'::jsonb
) ON CONFLICT (slug) DO NOTHING;

-- System roles
INSERT INTO roles (id, tenant_id, name, description, level, is_system) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'super_admin', 'Super Administrator - Full access', 10, true),
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'admin', 'Administrator', 8, true),
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'manager', 'Manager', 6, true),
    ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001', 'user', 'Regular User', 5, true),
    ('b0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'viewer', 'Read-only Viewer', 3, true)
ON CONFLICT DO NOTHING;

-- Permissions
INSERT INTO permissions (id, name, resource, action, description) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'books.create', 'books', 'create', 'Create books'),
    ('c0000000-0000-0000-0000-000000000002', 'books.read', 'books', 'read', 'Read books'),
    ('c0000000-0000-0000-0000-000000000003', 'books.update', 'books', 'update', 'Update books'),
    ('c0000000-0000-0000-0000-000000000004', 'books.delete', 'books', 'delete', 'Delete books'),
    ('c0000000-0000-0000-0000-000000000005', 'books.export', 'books', 'export', 'Export books'),
    ('c0000000-0000-0000-0000-000000000006', 'loans.create', 'loans', 'create', 'Create loans'),
    ('c0000000-0000-0000-0000-000000000007', 'loans.read', 'loans', 'read', 'Read loans'),
    ('c0000000-0000-0000-0000-000000000008', 'loans.update', 'loans', 'update', 'Update loans'),
    ('c0000000-0000-0000-0000-000000000009', 'loans.delete', 'loans', 'delete', 'Delete loans'),
    ('c0000000-0000-0000-0000-000000000010', 'users.read', 'users', 'read', 'Read users'),
    ('c0000000-0000-0000-0000-000000000011', 'users.manage', 'users', 'manage', 'Manage users'),
    ('c0000000-0000-0000-0000-000000000012', 'tenant.manage', 'tenant', 'manage', 'Manage tenant settings'),
    ('c0000000-0000-0000-0000-000000000013', 'tasks.create', 'tasks', 'create', 'Create tasks'),
    ('c0000000-0000-0000-0000-000000000014', 'tasks.read', 'tasks', 'read', 'Read tasks'),
    ('c0000000-0000-0000-0000-000000000015', 'tasks.update', 'tasks', 'update', 'Update tasks'),
    ('c0000000-0000-0000-0000-000000000016', 'tasks.delete', 'tasks', 'delete', 'Delete tasks')
ON CONFLICT DO NOTHING;

-- Role-Permission mappings
-- Super Admin: all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'b0000000-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT DO NOTHING;

-- Admin: all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'b0000000-0000-0000-0000-000000000002', id FROM permissions
ON CONFLICT DO NOTHING;

-- Manager: books and loans (CRUD) + tasks + users.read
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'b0000000-0000-0000-0000-000000000003', id FROM permissions
WHERE name IN ('books.create', 'books.read', 'books.update', 'books.delete',
               'loans.create', 'loans.read', 'loans.update', 'loans.delete',
               'tasks.create', 'tasks.read', 'tasks.update', 'tasks.delete',
               'users.read')
ON CONFLICT DO NOTHING;

-- User: books (create, read, update) + loans (create, read) + tasks
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'b0000000-0000-0000-0000-000000000004', id FROM permissions
WHERE name IN ('books.create', 'books.read', 'books.update',
               'loans.create', 'loans.read',
               'tasks.create', 'tasks.read', 'tasks.update', 'tasks.delete')
ON CONFLICT DO NOTHING;

-- Viewer: read-only
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'b0000000-0000-0000-0000-000000000005', id FROM permissions
WHERE name IN ('books.read', 'loans.read', 'tasks.read')
ON CONFLICT DO NOTHING;

-- Test users (passwordless - no password needed)
INSERT INTO users (id, tenant_id, email, full_name, is_active, role_id)
VALUES
    ('d0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'admin@apsferreira.com', 'Admin User', true, 'b0000000-0000-0000-0000-000000000002'),
    ('d0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'user@apsferreira.com', 'Regular User', true, 'b0000000-0000-0000-0000-000000000004')
ON CONFLICT DO NOTHING;

-- User-Role mappings
INSERT INTO user_roles (user_id, role_id) VALUES
    ('d0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000002'),
    ('d0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000004')
ON CONFLICT DO NOTHING;
