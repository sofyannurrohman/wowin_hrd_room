package repository

import (
	"context"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ParticipantRepository struct {
	db *pgxpool.Pool
}

func NewParticipantRepository(db *pgxpool.Pool) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) Create(ctx context.Context, p *domain.SessionParticipant) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO session_participants (id, session_id, user_id, token_id, joined_at, status)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		p.ID, p.SessionID, p.UserID, p.TokenID, time.Now(), p.Status)
	return err
}

func (r *ParticipantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SessionParticipant, error) {
	row := r.db.QueryRow(ctx,
		`SELECT sp.id, sp.session_id, sp.user_id, u.name, u.email, sp.token_id, sp.joined_at, sp.disconnected_at, sp.status
		 FROM session_participants sp JOIN users u ON sp.user_id = u.id WHERE sp.id=$1`, id)
	p := &domain.SessionParticipant{}
	return p, row.Scan(&p.ID, &p.SessionID, &p.UserID, &p.UserName, &p.UserEmail, &p.TokenID, &p.JoinedAt, &p.DisconnectedAt, &p.Status)
}

func (r *ParticipantRepository) FindByUserAndSession(ctx context.Context, userID, sessionID uuid.UUID) (*domain.SessionParticipant, error) {
	row := r.db.QueryRow(ctx,
		`SELECT sp.id, sp.session_id, sp.user_id, u.name, u.email, sp.token_id, sp.joined_at, sp.disconnected_at, sp.status
		 FROM session_participants sp JOIN users u ON sp.user_id = u.id
		 WHERE sp.user_id=$1 AND sp.session_id=$2`, userID, sessionID)
	p := &domain.SessionParticipant{}
	err := row.Scan(&p.ID, &p.SessionID, &p.UserID, &p.UserName, &p.UserEmail, &p.TokenID, &p.JoinedAt, &p.DisconnectedAt, &p.Status)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ParticipantRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.SessionParticipant, error) {
	rows, err := r.db.Query(ctx,
		`SELECT sp.id, sp.session_id, sp.user_id, u.name, u.email, sp.token_id, sp.joined_at, sp.disconnected_at, sp.status
		 FROM session_participants sp JOIN users u ON sp.user_id = u.id WHERE sp.session_id=$1 ORDER BY sp.joined_at ASC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.SessionParticipant
	for rows.Next() {
		var p domain.SessionParticipant
		if err := rows.Scan(&p.ID, &p.SessionID, &p.UserID, &p.UserName, &p.UserEmail, &p.TokenID, &p.JoinedAt, &p.DisconnectedAt, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *ParticipantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	now := time.Now()
	if status == "finished" || status == "disconnected" {
		_, err := r.db.Exec(ctx, `UPDATE session_participants SET status=$1, disconnected_at=$2 WHERE id=$3`, status, now, id)
		return err
	}
	_, err := r.db.Exec(ctx, `UPDATE session_participants SET status=$1 WHERE id=$2`, status, id)
	return err
}

type AnswerRepository struct {
	db *pgxpool.Pool
}

func NewAnswerRepository(db *pgxpool.Pool) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Upsert(ctx context.Context, a *domain.Answer) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO answers (id, participant_id, question_id, selected_option_id, text_answer, answered_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (participant_id, question_id)
		 DO UPDATE SET selected_option_id=$4, text_answer=$5, answered_at=$6`,
		a.ID, a.ParticipantID, a.QuestionID, a.SelectedOptionID, a.TextAnswer, time.Now())
	return err
}

func (r *AnswerRepository) ListByParticipant(ctx context.Context, participantID uuid.UUID) ([]domain.Answer, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, participant_id, question_id, selected_option_id, text_answer, is_correct, score, hr_notes, answered_at
		 FROM answers WHERE participant_id=$1`, participantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		var a domain.Answer
		if err := rows.Scan(&a.ID, &a.ParticipantID, &a.QuestionID, &a.SelectedOptionID, &a.TextAnswer, &a.IsCorrect, &a.Score, &a.HRNotes, &a.AnsweredAt); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (r *AnswerRepository) UpdateGrading(ctx context.Context, answerID uuid.UUID, isCorrect bool, score float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE answers SET is_correct=$1, score=$2 WHERE id=$3`, isCorrect, score, answerID)
	return err
}

func (r *AnswerRepository) UpdateHRReview(ctx context.Context, answerID uuid.UUID, score float64, notes string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE answers SET score=$1, hr_notes=$2, is_correct=true WHERE id=$3`, score, notes, answerID)
	return err
}

type ResultRepository struct {
	db *pgxpool.Pool
}

func NewResultRepository(db *pgxpool.Pool) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r *ResultRepository) Create(ctx context.Context, res *domain.Result) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO results (id, participant_id, total_score, grading_status, submitted_at)
		 VALUES ($1,$2,$3,$4,$5)`,
		res.ID, res.ParticipantID, res.TotalScore, res.GradingStatus, time.Now())
	return err
}

func (r *ResultRepository) FindByParticipant(ctx context.Context, participantID uuid.UUID) (*domain.Result, error) {
	row := r.db.QueryRow(ctx,
		`SELECT r.id, r.participant_id, u.name, u.email, r.total_score, r.grading_status, r.submitted_at
		 FROM results r JOIN session_participants sp ON r.participant_id = sp.id JOIN users u ON sp.user_id = u.id
		 WHERE r.participant_id=$1`, participantID)
	res := &domain.Result{}
	return res, row.Scan(&res.ID, &res.ParticipantID, &res.UserName, &res.UserEmail, &res.TotalScore, &res.GradingStatus, &res.SubmittedAt)
}

func (r *ResultRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.Result, error) {
	rows, err := r.db.Query(ctx,
		`SELECT r.id, r.participant_id, u.name, u.email, r.total_score, r.grading_status, r.submitted_at
		 FROM results r
		 JOIN session_participants sp ON r.participant_id = sp.id
		 JOIN users u ON sp.user_id = u.id
		 WHERE sp.session_id=$1 ORDER BY r.submitted_at DESC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.Result
	for rows.Next() {
		var res domain.Result
		if err := rows.Scan(&res.ID, &res.ParticipantID, &res.UserName, &res.UserEmail, &res.TotalScore, &res.GradingStatus, &res.SubmittedAt); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

func (r *ResultRepository) UpdateScore(ctx context.Context, id uuid.UUID, totalScore float64, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE results SET total_score=$1, grading_status=$2 WHERE id=$3`, totalScore, status, id)
	return err
}

type ViolationRepository struct {
	db *pgxpool.Pool
}

func NewViolationRepository(db *pgxpool.Pool) *ViolationRepository {
	return &ViolationRepository{db: db}
}

func (r *ViolationRepository) Create(ctx context.Context, v *domain.Violation) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO violations (id, participant_id, violation_type, detected_at, resolved)
		 VALUES ($1,$2,$3,$4,$5)`,
		v.ID, v.ParticipantID, v.ViolationType, time.Now(), false)
	return err
}

func (r *ViolationRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.Violation, error) {
	rows, err := r.db.Query(ctx,
		`SELECT v.id, v.participant_id, sp.session_id, u.name, v.violation_type, v.detected_at, v.resolved
		 FROM violations v
		 JOIN session_participants sp ON v.participant_id = sp.id
		 JOIN users u ON sp.user_id = u.id
		 WHERE sp.session_id=$1 ORDER BY v.detected_at DESC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var violations []domain.Violation
	for rows.Next() {
		var viol domain.Violation
		if err := rows.Scan(&viol.ID, &viol.ParticipantID, &viol.SessionID, &viol.UserName, &viol.ViolationType, &viol.DetectedAt, &viol.Resolved); err != nil {
			return nil, err
		}
		violations = append(violations, viol)
	}
	return violations, nil
}

func (r *ViolationRepository) ListByParticipant(ctx context.Context, participantID uuid.UUID) ([]domain.Violation, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, participant_id, '', '', violation_type, detected_at, resolved
		 FROM violations WHERE participant_id=$1 ORDER BY detected_at DESC`, participantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var violations []domain.Violation
	for rows.Next() {
		var v domain.Violation
		if err := rows.Scan(&v.ID, &v.ParticipantID, &v.SessionID, &v.UserName, &v.ViolationType, &v.DetectedAt, &v.Resolved); err != nil {
			return nil, err
		}
		violations = append(violations, v)
	}
	return violations, nil
}

type LogRepository struct {
	db *pgxpool.Pool
}

func NewLogRepository(db *pgxpool.Pool) *LogRepository {
	return &LogRepository{db: db}
}

func (r *LogRepository) Create(ctx context.Context, userID *uuid.UUID, action, detail string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO system_logs (user_id, action, detail) VALUES ($1,$2,$3)`, userID, action, detail)
	return err
}

func (r *LogRepository) List(ctx context.Context) ([]domain.SystemLog, error) {
	rows, err := r.db.Query(ctx,
		`SELECT sl.id, sl.user_id, COALESCE(u.name,'system'), sl.action, COALESCE(sl.detail,''), sl.created_at
		 FROM system_logs sl LEFT JOIN users u ON sl.user_id = u.id ORDER BY sl.created_at DESC LIMIT 500`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.SystemLog
	for rows.Next() {
		var l domain.SystemLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.UserName, &l.Action, &l.Detail, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
