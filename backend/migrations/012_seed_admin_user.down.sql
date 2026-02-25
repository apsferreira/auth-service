-- Migration 012 rollback: remove admin user apsferreira
DELETE FROM user_roles WHERE user_id = 'd0000000-0000-0000-0000-000000000003';
DELETE FROM users WHERE id = 'd0000000-0000-0000-0000-000000000003';
