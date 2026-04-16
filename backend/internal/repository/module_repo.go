package repository

import (
	"context"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ModuleRepository struct {
	db *pgxpool.Pool
}

func NewModuleRepository(db *pgxpool.Pool) *ModuleRepository {
	return &ModuleRepository{db: db}
}

func (r *ModuleRepository) Create(ctx context.Context, m *domain.Module) error {
	query := `INSERT INTO modules (id, name, description, memorization_content, memorization_duration, total_weight, created_by, created_at)
	          VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.Exec(ctx, query, m.ID, m.Name, m.Description, m.MemorizationContent, m.MemorizationDuration, m.TotalWeight, m.CreatedBy, time.Now())
	return err
}

func (r *ModuleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Module, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, description, memorization_content, memorization_duration, total_weight, created_by, created_at FROM modules WHERE id=$1`, id)
	m := &domain.Module{}
	err := row.Scan(&m.ID, &m.Name, &m.Description, &m.MemorizationContent, &m.MemorizationDuration, &m.TotalWeight, &m.CreatedBy, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *ModuleRepository) ListAll(ctx context.Context) ([]domain.Module, error) {
	query := `SELECT id, name, description, memorization_content, memorization_duration, total_weight, created_by, created_at FROM modules ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []domain.Module
	for rows.Next() {
		var m domain.Module
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.MemorizationContent, &m.MemorizationDuration, &m.TotalWeight, &m.CreatedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *ModuleRepository) Update(ctx context.Context, m *domain.Module) error {
	_, err := r.db.Exec(ctx,
		`UPDATE modules SET name=$1, description=$2, memorization_content=$3, memorization_duration=$4, total_weight=$5 WHERE id=$6`,
		m.Name, m.Description, m.MemorizationContent, m.MemorizationDuration, m.TotalWeight, m.ID)
	return err
}

func (r *ModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM modules WHERE id=$1`, id)
	return err
}

func (r *ModuleRepository) AddSessionModule(ctx context.Context, sm *domain.SessionModule) error {
	query := `INSERT INTO session_modules (session_id, module_id, sort_order) VALUES ($1,$2,$3) ON CONFLICT (session_id, module_id) DO UPDATE SET sort_order=$3`
	_, err := r.db.Exec(ctx, query, sm.SessionID, sm.ModuleID, sm.SortOrder)
	return err
}

func (r *ModuleRepository) DeleteSessionModules(ctx context.Context, sessionID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM session_modules WHERE session_id=$1`, sessionID)
	return err
}

func (r *ModuleRepository) ListBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.Module, error) {
	query := `
		SELECT m.id, m.name, m.description, m.memorization_content, m.memorization_duration, m.total_weight, m.created_by, m.created_at
		FROM modules m
		JOIN session_modules sm ON sm.module_id = m.id
		WHERE sm.session_id = $1
		ORDER BY sm.sort_order ASC
	`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []domain.Module
	for rows.Next() {
		var m domain.Module
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.MemorizationContent, &m.MemorizationDuration, &m.TotalWeight, &m.CreatedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}
