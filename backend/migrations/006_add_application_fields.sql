-- Adds new columns to users table for public job application
ALTER TABLE users ADD COLUMN cv_url VARCHAR(500);
ALTER TABLE users ADD COLUMN expected_salary VARCHAR(100);
