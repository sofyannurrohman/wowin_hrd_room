package usecase

import (
	"context"
	"hrd_room/backend/internal/repository"
)

type SettingsUseCase struct {
	settingsRepo *repository.SettingsRepository
}

func NewSettingsUseCase(sr *repository.SettingsRepository) *SettingsUseCase {
	return &SettingsUseCase{settingsRepo: sr}
}

func (uc *SettingsUseCase) GetSettings(ctx context.Context, key string) (interface{}, error) {
	var val interface{}
	err := uc.settingsRepo.Get(ctx, key, &val)
	return val, err
}

func (uc *SettingsUseCase) UpdateSettings(ctx context.Context, key string, value interface{}) error {
	return uc.settingsRepo.Update(ctx, key, value)
}
