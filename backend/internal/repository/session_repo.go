package repository

import (
	"context"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *domain.Session) error {
	query := `INSERT INTO sessions (id, created_by, name, schedule, duration_minutes, max_participants, randomize_questions, show_score, is_locked, created_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(ctx, query, s.ID, s.CreatedBy, s.Name, s.Schedule, s.DurationMinutes,
		s.MaxParticipants, s.RandomizeQuestions, s.ShowScore, s.IsLocked, time.Now())
	return err
}

func (r *SessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	query := `SELECT s.id, s.created_by, u.name, s.name, s.schedule, s.duration_minutes,
	          s.max_participants, s.randomize_questions, s.show_score, s.is_locked, s.created_at,
	          COUNT(sp.id) FILTER (WHERE sp.status = 'active') as active_participant_count,
	          COUNT(sp.id) FILTER (WHERE sp.status = 'finished') as submitted_participant_count,
	          COUNT(sp.id) as total_participant_count
	          FROM sessions s
	          LEFT JOIN users u ON s.created_by = u.id
	          LEFT JOIN session_participants sp ON sp.session_id = s.id
	          WHERE s.id = $1
	          GROUP BY s.id, u.name`
	row := r.db.QueryRow(ctx, query, id)
	s := &domain.Session{}
	err := row.Scan(&s.ID, &s.CreatedBy, &s.CreatedByName, &s.Name, &s.Schedule,
		&s.DurationMinutes, &s.MaxParticipants, &s.RandomizeQuestions,
		&s.ShowScore, &s.IsLocked, &s.CreatedAt, &s.ActiveParticipantCount,
		&s.SubmittedParticipantCount, &s.TotalParticipantCount)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SessionRepository) List(ctx context.Context, createdBy *uuid.UUID) ([]domain.Session, error) {
	query := `SELECT s.id, s.created_by, u.name, s.name, s.schedule, s.duration_minutes,
	          s.max_participants, s.randomize_questions, s.show_score, s.is_locked, s.created_at,
	          COUNT(sp.id) FILTER (WHERE sp.status = 'active') as active_participant_count,
	          COUNT(sp.id) FILTER (WHERE sp.status = 'finished') as submitted_participant_count,
	          COUNT(sp.id) as total_participant_count
	          FROM sessions s
	          LEFT JOIN users u ON s.created_by = u.id
	          LEFT JOIN session_participants sp ON sp.session_id = s.id`
	args := []interface{}{}
	if createdBy != nil {
		query += ` WHERE s.created_by = $1`
		args = append(args, *createdBy)
	}
	query += ` GROUP BY s.id, u.name ORDER BY s.created_at DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.Session
	for rows.Next() {
		var s domain.Session
		if err := rows.Scan(&s.ID, &s.CreatedBy, &s.CreatedByName, &s.Name, &s.Schedule,
			&s.DurationMinutes, &s.MaxParticipants, &s.RandomizeQuestions,
			&s.ShowScore, &s.IsLocked, &s.CreatedAt, &s.ActiveParticipantCount,
			&s.SubmittedParticipantCount, &s.TotalParticipantCount); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, s *domain.Session) error {
	_, err := r.db.Exec(ctx, `UPDATE sessions SET name=$1, schedule=$2, duration_minutes=$3,
	max_participants=$4, randomize_questions=$5, show_score=$6, is_locked=$7 WHERE id=$8`,
		s.Name, s.Schedule, s.DurationMinutes, s.MaxParticipants,
		s.RandomizeQuestions, s.ShowScore, s.IsLocked, s.ID)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE id=$1`, id)
	return err
}

func (r *SessionRepository) Lock(ctx context.Context, id uuid.UUID, locked bool) error {
	_, err := r.db.Exec(ctx, `UPDATE sessions SET is_locked=$1 WHERE id=$2`, locked, id)
	return err
}

func (r *SessionRepository) CountActiveParticipants(ctx context.Context, sessionID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM session_participants WHERE session_id=$1 AND status='active'`, sessionID).Scan(&count)
	return count, err
}
