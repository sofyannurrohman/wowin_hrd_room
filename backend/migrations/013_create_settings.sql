-- Migration 013: Global Settings
CREATE TABLE settings (
    key VARCHAR(100) PRIMARY KEY,
    value JSONB NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed default settings
INSERT INTO settings (key, value) VALUES 
('notification_settings', '{"sessionCompleted": true, "manualGrading": true, "candidateShortlisted": false}'),
('system_settings', '{"strictMode": false, "dataRetention": "Keep Forever"}');
