-- Unlink permissions from services
UPDATE permissions SET service_id = NULL;

-- Remove seeded services
DELETE FROM services WHERE id IN (
  'e0000000-0000-0000-0000-000000000001',
  'e0000000-0000-0000-0000-000000000002',
  'e0000000-0000-0000-0000-000000000003'
);
