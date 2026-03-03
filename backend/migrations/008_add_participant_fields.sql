-- Migration 008: Add participant fields to users

ALTER TABLE users 
ADD COLUMN age INT,
ADD COLUMN applied_position VARCHAR(100);
