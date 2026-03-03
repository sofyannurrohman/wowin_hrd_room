package usecase

import (
	"context"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"

	"github.com/google/uuid"
)

type JobPositionUseCase struct {
	repo *repository.JobPositionRepository
}

func NewJobPositionUseCase(repo *repository.JobPositionRepository) *JobPositionUseCase {
	return &JobPositionUseCase{repo: repo}
}

func (uc *JobPositionUseCase) GetActivePositions(ctx context.Context) ([]*domain.JobPosition, error) {
	return uc.repo.GetActive(ctx)
}

func (uc *JobPositionUseCase) GetAllPositions(ctx context.Context) ([]*domain.JobPosition, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *JobPositionUseCase) GetPosition(ctx context.Context, id uuid.UUID) (*domain.JobPosition, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *JobPositionUseCase) CreatePosition(ctx context.Context, p *domain.JobPosition) error {
	return uc.repo.Create(ctx, p)
}

func (uc *JobPositionUseCase) UpdatePosition(ctx context.Context, id uuid.UUID, p *domain.JobPosition) error {
	return uc.repo.Update(ctx, id, p)
}

func (uc *JobPositionUseCase) DeletePosition(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
