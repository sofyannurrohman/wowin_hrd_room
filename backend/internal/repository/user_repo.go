package repository

import (
	"context"
	"fmt"
	"time"

	"hrd_room/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO users (id, role_id, name, email, password_hash, age, applied_position, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, u.ID, u.RoleID, u.Name, u.Email, u.PasswordHash, u.Age, u.AppliedPosition, time.Now())
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT u.id, u.role_id, r.name, u.name, u.email, u.password_hash, u.age, u.applied_position, u.created_at
	          FROM users u JOIN roles r ON u.role_id = r.id WHERE u.email = $1`
	row := r.db.QueryRow(ctx, query, email)
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Email, &u.PasswordHash, &u.Age, &u.AppliedPosition, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT u.id, u.role_id, r.name, u.name, u.email, u.password_hash, u.age, u.applied_position, u.created_at
	          FROM users u JOIN roles r ON u.role_id = r.id WHERE u.id = $1`
	row := r.db.QueryRow(ctx, query, id)
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Email, &u.PasswordHash, &u.Age, &u.AppliedPosition, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	query := `SELECT u.id, u.role_id, r.name, u.name, u.email, u.age, u.applied_position, u.created_at
	          FROM users u JOIN roles r ON u.role_id = r.id ORDER BY u.created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Email, &u.Age, &u.AppliedPosition, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET name=$1, email=$2, role_id=$3, age=$4, applied_position=$5 WHERE id=$6`,
		u.Name, u.Email, u.RoleID, u.Age, u.AppliedPosition, u.ID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET password_hash=$1 WHERE id=$2`, hash, id)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
