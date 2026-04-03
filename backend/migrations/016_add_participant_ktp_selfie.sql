-- Migration 016: Add KTP Selfie URL to session_participants

ALTER TABLE session_participants 
ADD COLUMN ktp_selfie_url VARCHAR(500);
