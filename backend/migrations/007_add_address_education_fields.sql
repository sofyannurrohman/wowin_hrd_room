-- Adds new columns to users table for public job application
ALTER TABLE users ADD COLUMN address VARCHAR(500);
ALTER TABLE users ADD COLUMN last_education VARCHAR(100);
