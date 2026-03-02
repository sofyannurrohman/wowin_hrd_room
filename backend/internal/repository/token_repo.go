package repository

import (
	"context"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, t *domain.ExamToken) error {
	query := `INSERT INTO exam_tokens (id, session_id, token_hash, token_type, max_usage, used_count, expires_at, is_active, bound_email, created_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(ctx, query, t.ID, t.SessionID, t.TokenHash, t.TokenType,
		t.MaxUsage, t.UsedCount, t.ExpiresAt, t.IsActive, t.BoundEmail, time.Now())
	return err
}

func (r *TokenRepository) FindByHash(ctx context.Context, hash string) (*domain.ExamToken, error) {
	query := `SELECT id, session_id, token_hash, token_type, max_usage, used_count, expires_at, is_active, bound_email, created_at
	          FROM exam_tokens WHERE token_hash=$1`
	row := r.db.QueryRow(ctx, query, hash)
	t := &domain.ExamToken{}
	err := row.Scan(&t.ID, &t.SessionID, &t.TokenHash, &t.TokenType, &t.MaxUsage,
		&t.UsedCount, &t.ExpiresAt, &t.IsActive, &t.BoundEmail, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TokenRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.ExamToken, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, session_id, token_hash, token_type, max_usage, used_count, expires_at, is_active, bound_email, created_at
		 FROM exam_tokens WHERE session_id=$1 ORDER BY created_at DESC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.ExamToken
	for rows.Next() {
		var t domain.ExamToken
		if err := rows.Scan(&t.ID, &t.SessionID, &t.TokenHash, &t.TokenType, &t.MaxUsage,
			&t.UsedCount, &t.ExpiresAt, &t.IsActive, &t.BoundEmail, &t.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (r *TokenRepository) IncrementUsage(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE exam_tokens SET used_count = used_count + 1 WHERE id=$1`, id)
	return err
}

func (r *TokenRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE exam_tokens SET is_active=false WHERE id=$1`, id)
	return err
}
