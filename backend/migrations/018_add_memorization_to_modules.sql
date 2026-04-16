-- Migration 018: Add memorization fields to modules

ALTER TABLE modules ADD COLUMN memorization_content TEXT;
ALTER TABLE modules ADD COLUMN memorization_duration INT DEFAULT 0; -- duration in seconds
