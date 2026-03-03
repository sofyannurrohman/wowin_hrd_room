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
