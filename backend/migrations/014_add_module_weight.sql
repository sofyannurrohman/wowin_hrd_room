-- Migration 014: Add total_weight to modules

ALTER TABLE modules ADD COLUMN IF NOT EXISTS total_weight DECIMAL(5,2) DEFAULT 100.00;
