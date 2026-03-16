package repository

import (
	"context"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobPositionRepository struct {
	db *pgxpool.Pool
}

func NewJobPositionRepository(db *pgxpool.Pool) *JobPositionRepository {
	return &JobPositionRepository{db: db}
}

func (r *JobPositionRepository) GetActive(ctx context.Context) ([]*domain.JobPosition, error) {
	query := `SELECT id, name, description, requirements, is_active, created_at, updated_at FROM job_positions WHERE is_active = true ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*domain.JobPosition
	for rows.Next() {
		p := &domain.JobPosition{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Requirements, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *JobPositionRepository) GetAll(ctx context.Context) ([]*domain.JobPosition, error) {
	query := `SELECT id, name, description, requirements, is_active, created_at, updated_at FROM job_positions ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*domain.JobPosition
	for rows.Next() {
		p := &domain.JobPosition{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Requirements, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (r *JobPositionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.JobPosition, error) {
	query := `SELECT id, name, description, requirements, is_active, created_at, updated_at FROM job_positions WHERE id = $1`
	p := &domain.JobPosition{}
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Requirements, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *JobPositionRepository) Create(ctx context.Context, p *domain.JobPosition) error {
	query := `INSERT INTO job_positions (name, description, requirements, is_active) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, p.Name, p.Description, p.Requirements, p.IsActive).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *JobPositionRepository) Update(ctx context.Context, id uuid.UUID, p *domain.JobPosition) error {
	query := `UPDATE job_positions SET name = $1, description = $2, requirements = $3, is_active = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, p.Name, p.Description, p.Requirements, p.IsActive, id)
	return err
}

func (r *JobPositionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM job_positions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
