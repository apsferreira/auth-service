-- Migration 010: Add performance indexes
-- Composite and additional indexes for multi-tenant queries and auth operations

-- Users: composite index for tenant + email lookup (common auth path)
CREATE INDEX IF NOT EXISTS idx_users_tenant_email ON users(tenant_id, email);

-- Users: composite index for listing active users per tenant
CREATE INDEX IF NOT EXISTS idx_users_tenant_active ON users(tenant_id, is_active);

-- Roles: filter system vs custom roles
CREATE INDEX IF NOT EXISTS idx_roles_is_system ON roles(is_system);

-- OTP: composite index for verification (email + expiry checked together)
CREATE INDEX IF NOT EXISTS idx_otp_codes_email_expires ON otp_codes(email, expires_at);

-- Refresh tokens: for cleanup queries on revoked tokens
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_revoked_at ON refresh_tokens(revoked_at);

-- Services: composite index for listing active services per tenant
CREATE INDEX IF NOT EXISTS idx_services_tenant_active ON services(tenant_id, is_active);

-- Permissions: index on resource + action for authorization checks
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);

-- Role permissions: reverse lookup (which roles have a specific permission)
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_role ON role_permissions(permission_id, role_id);
