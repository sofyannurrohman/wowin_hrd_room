-- Migration 010: Create job positions

CREATE TABLE job_positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_job_positions
BEFORE UPDATE ON job_positions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Insert default positions
INSERT INTO job_positions (name, is_active) VALUES
('Software Engineer', true),
('Product Manager', true),
('Data Analyst', true),
('UI/UX Designer', true),
('Marketing Specialist', true),
('Other', true);
