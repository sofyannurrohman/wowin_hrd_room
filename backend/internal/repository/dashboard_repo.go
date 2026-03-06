package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

type DashboardStats struct {
	ActiveSessions         int `json:"active_sessions"`
	TotalParticipantsToday int `json:"total_participants_today"`
	RecentViolations       int `json:"recent_violations"`
}

func (r *DashboardRepository) GetStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	err := r.db.QueryRow(ctx, `SELECT count(*) FROM sessions WHERE is_locked=false`).Scan(&stats.ActiveSessions)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `SELECT count(*) FROM session_participants WHERE joined_at >= $1`, todayStart).Scan(&stats.TotalParticipantsToday)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `SELECT count(*) FROM violations WHERE detected_at >= $1`, todayStart).Scan(&stats.RecentViolations)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

type ParticipantStat struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	LastSessionName string    `json:"last_session_name"`
	LastSessionDate time.Time `json:"last_session_date"`
	AvgScore        float64   `json:"avg_score"`
	Status          string    `json:"status"`
}

func (r *DashboardRepository) ListParticipants(ctx context.Context) ([]ParticipantStat, error) {
	query := `
        SELECT CAST(u.id AS TEXT), u.name, u.email, 
               COALESCE(s.name, ''), COALESCE(sp.joined_at, '1970-01-01'::timestamp),
               COALESCE((SELECT avg(total_score) FROM results r JOIN session_participants sp2 ON r.participant_id = sp2.id WHERE sp2.user_id = u.id), 0) as avg_score
        FROM users u
        LEFT JOIN (
            SELECT user_id, max(joined_at) as last_join FROM session_participants GROUP BY user_id
        ) latest_sp ON latest_sp.user_id = u.id
        LEFT JOIN session_participants sp ON sp.user_id = u.id AND sp.joined_at = latest_sp.last_join
        LEFT JOIN sessions s ON sp.session_id = s.id
        WHERE u.role_id = 3
        ORDER BY u.created_at DESC
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ParticipantStat
	for rows.Next() {
		var p ParticipantStat
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.LastSessionName, &p.LastSessionDate, &p.AvgScore); err != nil {
			return nil, err
		}
		p.Status = "Active"
		list = append(list, p)
	}
	return list, nil
}

func (r *DashboardRepository) GetParticipantDetail(ctx context.Context, userID string) (map[string]interface{}, error) {
	var p struct {
		ID              string  `json:"id"`
		Name            string  `json:"name"`
		Email           string  `json:"email"`
		Age             *int    `json:"age"`
		AppliedPosition *string `json:"applied_position"`
		Address         *string `json:"address"`
		LastEducation   *string `json:"last_education"`
		WhatsappNumber  *string `json:"whatsapp_number"`
		CvURL           *string `json:"cv_url"`
		ExpectedSalary  *string `json:"expected_salary"`
	}

	err := r.db.QueryRow(ctx, `
		SELECT CAST(id AS TEXT), name, email, age, applied_position, address, last_education, whatsapp_number, cv_url, expected_salary
		FROM users WHERE id = $1 AND role_id = 3
	`, userID).Scan(&p.ID, &p.Name, &p.Email, &p.Age, &p.AppliedPosition, &p.Address, &p.LastEducation, &p.WhatsappNumber, &p.CvURL, &p.ExpectedSalary)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT 
			s.name as session_name,
			s.schedule,
			sp.status as participant_status,
			r.total_score,
			CAST(r.id AS TEXT) as result_id,
            CAST(s.id AS TEXT) as session_id,
			CAST(sp.id AS TEXT) as participant_id
		FROM session_participants sp
		JOIN sessions s ON sp.session_id = s.id
		LEFT JOIN results r ON sp.id = r.participant_id
		WHERE sp.user_id = $1
		ORDER BY s.schedule DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var sessionName string
		var schedule time.Time
		var status string
		var score *float64
		var resultID *string
		var sessionID string
		var participantID string

		if err := rows.Scan(&sessionName, &schedule, &status, &score, &resultID, &sessionID, &participantID); err != nil {
			return nil, err
		}

		history = append(history, map[string]interface{}{
			"session_name":   sessionName,
			"schedule":       schedule,
			"status":         status,
			"score":          score,
			"result_id":      resultID,
			"session_id":     sessionID,
			"participant_id": participantID,
		})
	}
	if history == nil {
		history = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"profile": p,
		"history": history,
	}, nil
}
