package repository

import (
	"context"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewQuestionRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(ctx context.Context, q *domain.Question) error {
	query := `INSERT INTO questions (id, session_id, type, content, image_url, weight, requires_manual_review, created_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.Exec(ctx, query, q.ID, q.SessionID, q.Type, q.Content, q.ImageURL,
		q.Weight, q.RequiresManualReview, time.Now())
	return err
}

func (r *QuestionRepository) CreateOption(ctx context.Context, opt *domain.QuestionOption) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO question_options (id, question_id, content, is_correct) VALUES ($1,$2,$3,$4)`,
		opt.ID, opt.QuestionID, opt.Content, opt.IsCorrect)
	return err
}

func (r *QuestionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, session_id, type, content, image_url, weight, requires_manual_review, created_at FROM questions WHERE id=$1`, id)
	q := &domain.Question{}
	err := row.Scan(&q.ID, &q.SessionID, &q.Type, &q.Content, &q.ImageURL, &q.Weight, &q.RequiresManualReview, &q.CreatedAt)
	if err != nil {
		return nil, err
	}
	q.Options, _ = r.GetOptions(ctx, id)
	return q, nil
}

func (r *QuestionRepository) ListBySession(ctx context.Context, sessionID uuid.UUID, randomize bool) ([]domain.Question, error) {
	query := `SELECT id, session_id, type, content, image_url, weight, requires_manual_review, created_at FROM questions WHERE session_id=$1`
	if randomize {
		query += ` ORDER BY RANDOM()`
	} else {
		query += ` ORDER BY created_at ASC`
	}
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []domain.Question
	for rows.Next() {
		var q domain.Question
		if err := rows.Scan(&q.ID, &q.SessionID, &q.Type, &q.Content, &q.ImageURL, &q.Weight, &q.RequiresManualReview, &q.CreatedAt); err != nil {
			return nil, err
		}
		q.Options, _ = r.GetOptions(ctx, q.ID)
		questions = append(questions, q)
	}
	return questions, nil
}

func (r *QuestionRepository) GetOptions(ctx context.Context, questionID uuid.UUID) ([]domain.QuestionOption, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, question_id, content, is_correct FROM question_options WHERE question_id=$1 ORDER BY id`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var opts []domain.QuestionOption
	for rows.Next() {
		var o domain.QuestionOption
		if err := rows.Scan(&o.ID, &o.QuestionID, &o.Content, &o.IsCorrect); err != nil {
			return nil, err
		}
		opts = append(opts, o)
	}
	return opts, nil
}

func (r *QuestionRepository) Update(ctx context.Context, q *domain.Question) error {
	_, err := r.db.Exec(ctx,
		`UPDATE questions SET type=$1, content=$2, image_url=$3, weight=$4, requires_manual_review=$5 WHERE id=$6`,
		q.Type, q.Content, q.ImageURL, q.Weight, q.RequiresManualReview, q.ID)
	return err
}

func (r *QuestionRepository) DeleteOptions(ctx context.Context, questionID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM question_options WHERE question_id=$1`, questionID)
	return err
}

func (r *QuestionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM questions WHERE id=$1`, id)
	return err
}

func (r *QuestionRepository) UpdateImageURL(ctx context.Context, id uuid.UUID, imageURL string) error {
	_, err := r.db.Exec(ctx, `UPDATE questions SET image_url=$1 WHERE id=$2`, imageURL, id)
	return err
}
