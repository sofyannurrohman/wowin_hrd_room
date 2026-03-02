-- Seed data for development
-- Run AFTER 001_init.sql

-- Default HR user: hr@hrdroom.com / password123
-- Default Admin: admin@hrdroom.com / admin123

INSERT INTO users (id, role_id, name, email, password_hash) VALUES
  ('11111111-1111-1111-1111-111111111111', 1, 'Super Admin', 'admin@hrdroom.com',
   '$2a$10$YourHashHereWillBeReplacedByScript'),
  ('22222222-2222-2222-2222-222222222222', 2, 'HR Team', 'hr@hrdroom.com',
   '$2a$10$YourHashHereWillBeReplacedByScript');

-- The app seed script (scripts/seed.go) will insert users with proper bcrypt hashes
