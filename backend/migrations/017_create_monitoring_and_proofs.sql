-- Add proof_url column to violations table
ALTER TABLE violations ADD COLUMN proof_url TEXT;

-- Create monitoring_photos table
CREATE TABLE IF NOT EXISTS monitoring_photos (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    participant_id UUID NOT NULL REFERENCES session_participants(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indices for performance
CREATE INDEX idx_monitoring_photos_participant ON monitoring_photos(participant_id);
CREATE INDEX idx_monitoring_photos_session ON monitoring_photos(session_id);
CREATE INDEX idx_monitoring_photos_created_at ON monitoring_photos(created_at);
