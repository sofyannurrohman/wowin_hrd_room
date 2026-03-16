package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepository struct {
	db *pgxpool.Pool
}

func NewSettingsRepository(db *pgxpool.Pool) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context, key string, target interface{}) error {
	var val []byte
	err := r.db.QueryRow(ctx, `SELECT value FROM settings WHERE key = $1`, key).Scan(&val)
	if err != nil {
		return err
	}
	return json.Unmarshal(val, target)
}

func (r *SettingsRepository) Update(ctx context.Context, key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `INSERT INTO settings (key, value, updated_at) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = $3`, 
		key, val, time.Now())
	return err
}
