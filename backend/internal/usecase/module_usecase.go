package usecase

import (
	"context"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"

	"github.com/google/uuid"
)

type ModuleUseCase struct {
	repo *repository.ModuleRepository
}

func NewModuleUseCase(repo *repository.ModuleRepository) *ModuleUseCase {
	return &ModuleUseCase{repo: repo}
}

func (uc *ModuleUseCase) CreateModule(ctx context.Context, m *domain.Module) error {
	m.ID = uuid.New()
	return uc.repo.Create(ctx, m)
}

func (uc *ModuleUseCase) UpdateModule(ctx context.Context, id uuid.UUID, m *domain.Module) error {
	current, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	current.Name = m.Name
	current.Description = m.Description
	return uc.repo.Update(ctx, current)
}

func (uc *ModuleUseCase) GetModule(ctx context.Context, id uuid.UUID) (*domain.Module, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *ModuleUseCase) ListModules(ctx context.Context) ([]domain.Module, error) {
	return uc.repo.ListAll(ctx)
}

func (uc *ModuleUseCase) DeleteModule(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
