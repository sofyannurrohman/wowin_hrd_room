package repository

import (
	"context"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MonitoringRepository struct {
	db *pgxpool.Pool
}

func NewMonitoringRepository(db *pgxpool.Pool) *MonitoringRepository {
	return &MonitoringRepository{db: db}
}

func (r *MonitoringRepository) Create(ctx context.Context, p *domain.MonitoringPhoto) error {
	query := `INSERT INTO monitoring_photos (id, session_id, participant_id, photo_url, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, p.ID, p.SessionID, p.ParticipantID, p.PhotoURL, p.CreatedAt)
	return err
}

func (r *MonitoringRepository) ListByParticipant(ctx context.Context, participantID uuid.UUID) ([]domain.MonitoringPhoto, error) {
	query := `SELECT id, session_id, participant_id, photo_url, created_at
	          FROM monitoring_photos WHERE participant_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, participantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []domain.MonitoringPhoto
	for rows.Next() {
		var p domain.MonitoringPhoto
		if err := rows.Scan(&p.ID, &p.SessionID, &p.ParticipantID, &p.PhotoURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *MonitoringRepository) GetOldPhotos(ctx context.Context, days int) ([]domain.MonitoringPhoto, error) {
	threshold := time.Now().AddDate(0, 0, -days)
	query := `SELECT id, session_id, participant_id, photo_url, created_at
	          FROM monitoring_photos WHERE created_at < $1`
	rows, err := r.db.Query(ctx, query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []domain.MonitoringPhoto
	for rows.Next() {
		var p domain.MonitoringPhoto
		if err := rows.Scan(&p.ID, &p.SessionID, &p.ParticipantID, &p.PhotoURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *MonitoringRepository) GetOldViolations(ctx context.Context, days int) ([]domain.Violation, error) {
	threshold := time.Now().AddDate(0, 0, -days)
	query := `SELECT id, participant_id, violation_type, detected_at, proof_url
	          FROM violations WHERE detected_at < $1 AND proof_url IS NOT NULL`
	rows, err := r.db.Query(ctx, query, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viols []domain.Violation
	for rows.Next() {
		var v domain.Violation
		if err := rows.Scan(&v.ID, &v.ParticipantID, &v.ViolationType, &v.DetectedAt, &v.ProofURL); err != nil {
			return nil, err
		}
		viols = append(viols, v)
	}
	return viols, nil
}

func (r *MonitoringRepository) DeletePhotos(ctx context.Context, ids []uuid.UUID) error {
	query := `DELETE FROM monitoring_photos WHERE id = ANY($1)`
	_, err := r.db.Exec(ctx, query, ids)
	return err
}

func (r *MonitoringRepository) ClearProofURLs(ctx context.Context, ids []uuid.UUID) error {
	query := `UPDATE violations SET proof_url = NULL WHERE id = ANY($1)`
	_, err := r.db.Exec(ctx, query, ids)
	return err
}
