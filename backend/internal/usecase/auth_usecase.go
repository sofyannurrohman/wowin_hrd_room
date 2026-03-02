package usecase

import (
	"context"
	"errors"
	"time"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"
	jwtpkg "hrd_room/backend/pkg/jwt"

	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo   *repository.UserRepository
	jwtManager *jwtpkg.Manager
}

func NewAuthUseCase(userRepo *repository.UserRepository, jwtManager *jwtpkg.Manager) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, jwtManager: jwtManager}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (uc *AuthUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !repository.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := uc.jwtManager.Generate(user.ID, user.Email, user.RoleName)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, User: *user}, nil
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   int    `json:"role_id"`
}

func (uc *AuthUseCase) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	hash, err := repository.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	roleID := req.RoleID
	if roleID == 0 {
		roleID = 3 // Peserta by default
	}

	user := &domain.User{
		ID:           uuid.New(),
		RoleID:       roleID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
