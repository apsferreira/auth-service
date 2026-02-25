-- Migration 012: Seed admin user apsferreira with password authentication
-- bcrypt hash of "Quixabeira@1" with cost 12

INSERT INTO users (id, tenant_id, email, full_name, username, password_hash, is_active, role_id)
VALUES (
    'd0000000-0000-0000-0000-000000000003',
    'a0000000-0000-0000-0000-000000000001',
    'apsf88@gmail.com',
    'Antonio Pedro Ferreira',
    'apsferreira',
    '$2a$12$kTHb3r7GdKRunW7MOgJE4uWPPI10H219nqHz01ZiJ4A0Lpu50NCWG',
    true,
    'b0000000-0000-0000-0000-000000000001'
)
ON CONFLICT (tenant_id, email) DO UPDATE SET
    username = EXCLUDED.username,
    password_hash = EXCLUDED.password_hash,
    full_name = EXCLUDED.full_name,
    role_id = EXCLUDED.role_id,
    is_active = true;

-- Ensure super_admin role is assigned in user_roles
INSERT INTO user_roles (user_id, role_id)
VALUES ('d0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;
