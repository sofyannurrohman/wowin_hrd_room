-- Migration 001: Create roles and users

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO roles (name) VALUES ('Super Admin'), ('HR'), ('Peserta');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id INT REFERENCES roles(id),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migration 002: Sessions and tokens

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_by UUID REFERENCES users(id),
    name VARCHAR(150) NOT NULL,
    schedule TIMESTAMP NOT NULL,
    duration_minutes INT NOT NULL,
    max_participants INT DEFAULT 20,
    randomize_questions BOOLEAN DEFAULT FALSE,
    show_score BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE exam_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    token_type VARCHAR(20) DEFAULT 'single-use',
    max_usage INT DEFAULT 1,
    used_count INT DEFAULT 0,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    bound_email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migration 003: Participants

CREATE TABLE session_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),
    token_id UUID REFERENCES exam_tokens(id),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    disconnected_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'active'
);

-- Migration 004: Questions

CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    image_url VARCHAR(500),
    weight DECIMAL(5,2) DEFAULT 1.00,
    requires_manual_review BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE question_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE
);

-- Migration 005: Answers and results

CREATE TABLE answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID REFERENCES session_participants(id) ON DELETE CASCADE,
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id UUID REFERENCES question_options(id),
    text_answer TEXT,
    is_correct BOOLEAN,
    score DECIMAL(5,2),
    hr_notes TEXT,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID REFERENCES session_participants(id) ON DELETE CASCADE,
    total_score DECIMAL(5,2) DEFAULT 0.00,
    grading_status VARCHAR(50) DEFAULT 'completed',
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migration 006: Violations and camera logs

CREATE TABLE violations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID REFERENCES session_participants(id) ON DELETE CASCADE,
    violation_type VARCHAR(100) NOT NULL,
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved BOOLEAN DEFAULT FALSE
);

CREATE TABLE camera_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    participant_id UUID REFERENCES session_participants(id) ON DELETE CASCADE,
    log_data JSONB,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migration 007: System logs

CREATE TABLE system_logs (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    detail TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_sessions_created_by ON sessions(created_by);
CREATE INDEX idx_exam_tokens_session_id ON exam_tokens(session_id);
CREATE INDEX idx_exam_tokens_hash ON exam_tokens(token_hash);
CREATE INDEX idx_participants_session ON session_participants(session_id);
CREATE INDEX idx_participants_user ON session_participants(user_id);
CREATE INDEX idx_questions_session ON questions(session_id);
CREATE INDEX idx_answers_participant ON answers(participant_id);
CREATE INDEX idx_violations_participant ON violations(participant_id);
CREATE INDEX idx_results_participant ON results(participant_id);
