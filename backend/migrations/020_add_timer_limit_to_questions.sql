-- Migration 020: Add timer_limit to questions
ALTER TABLE questions ADD COLUMN timer_limit INT DEFAULT 0;
COMMENT ON COLUMN questions.timer_limit IS 'Time limit for specific question in seconds (useful for typing tests)';
