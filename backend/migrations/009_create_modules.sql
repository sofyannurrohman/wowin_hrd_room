-- Migration 009: Create modules

CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE session_modules (
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    sort_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (session_id, module_id)
);

-- Migrate existing questions
-- Creates a dummy module just so we don't lose them, or we can just drop session_id
-- We'll just add module_id. If a question was tied to a session, maybe we could create a module for that session.
-- For now, let's just add module_id and drop session_id.

ALTER TABLE questions ADD COLUMN module_id UUID REFERENCES modules(id) ON DELETE CASCADE;
ALTER TABLE questions DROP COLUMN IF EXISTS session_id CASCADE;
