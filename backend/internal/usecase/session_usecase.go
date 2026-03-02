package usecase

import (
	"context"
	"errors"
	"time"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"
	"hrd_room/backend/pkg/hasher"

	"github.com/google/uuid"
)

type SessionUseCase struct {
	sessionRepo *repository.SessionRepository
	tokenRepo   *repository.TokenRepository
	logRepo     *repository.LogRepository
}

func NewSessionUseCase(sr *repository.SessionRepository, tr *repository.TokenRepository, lr *repository.LogRepository) *SessionUseCase {
	return &SessionUseCase{sessionRepo: sr, tokenRepo: tr, logRepo: lr}
}

type CreateSessionRequest struct {
	Name               string    `json:"name" binding:"required"`
	Schedule           time.Time `json:"schedule" binding:"required"`
	DurationMinutes    int       `json:"duration_minutes" binding:"required,min=1"`
	MaxParticipants    int       `json:"max_participants"`
	RandomizeQuestions bool      `json:"randomize_questions"`
	ShowScore          bool      `json:"show_score"`
}

func (uc *SessionUseCase) Create(ctx context.Context, createdBy uuid.UUID, req CreateSessionRequest) (*domain.Session, error) {
	maxP := req.MaxParticipants
	if maxP == 0 || maxP > 20 {
		maxP = 20
	}

	s := &domain.Session{
		ID:                 uuid.New(),
		CreatedBy:          createdBy,
		Name:               req.Name,
		Schedule:           req.Schedule,
		DurationMinutes:    req.DurationMinutes,
		MaxParticipants:    maxP,
		RandomizeQuestions: req.RandomizeQuestions,
		ShowScore:          req.ShowScore,
	}

	if err := uc.sessionRepo.Create(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (uc *SessionUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	return uc.sessionRepo.FindByID(ctx, id)
}

func (uc *SessionUseCase) List(ctx context.Context, role string, userID uuid.UUID) ([]domain.Session, error) {
	if role == domain.RoleSuperAdmin {
		return uc.sessionRepo.List(ctx, nil)
	}
	return uc.sessionRepo.List(ctx, &userID)
}

func (uc *SessionUseCase) Update(ctx context.Context, id uuid.UUID, req CreateSessionRequest) (*domain.Session, error) {
	s, err := uc.sessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.Name = req.Name
	s.Schedule = req.Schedule
	s.DurationMinutes = req.DurationMinutes
	s.MaxParticipants = req.MaxParticipants
	s.RandomizeQuestions = req.RandomizeQuestions
	s.ShowScore = req.ShowScore
	return s, uc.sessionRepo.Update(ctx, s)
}

func (uc *SessionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.sessionRepo.Delete(ctx, id)
}

// Token generation

type GenerateTokenRequest struct {
	Count      int        `json:"count" binding:"required,min=1,max=20"`
	TokenType  string     `json:"token_type"` // single-use, multi-use
	MaxUsage   int        `json:"max_usage"`
	ExpiresAt  *time.Time `json:"expires_at"`
	BoundEmail *string    `json:"bound_email"`
}

type GeneratedToken struct {
	ID         uuid.UUID  `json:"id"`
	PlainToken string     `json:"token"`
	TokenType  string     `json:"token_type"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

func (uc *SessionUseCase) GenerateTokens(ctx context.Context, sessionID uuid.UUID, req GenerateTokenRequest) ([]GeneratedToken, error) {
	if req.TokenType == "" {
		req.TokenType = "single-use"
	}
	maxUsage := req.MaxUsage
	if req.TokenType == "single-use" {
		maxUsage = 1
	}
	if maxUsage == 0 {
		maxUsage = 1
	}

	var results []GeneratedToken
	for i := 0; i < req.Count; i++ {
		plain, err := hasher.GenerateToken(48)
		if err != nil {
			return nil, err
		}

		hash := hasher.HashToken(plain)
		t := &domain.ExamToken{
			ID:         uuid.New(),
			SessionID:  sessionID,
			TokenHash:  hash,
			TokenType:  req.TokenType,
			MaxUsage:   maxUsage,
			UsedCount:  0,
			ExpiresAt:  req.ExpiresAt,
			IsActive:   true,
			BoundEmail: req.BoundEmail,
		}

		if err := uc.tokenRepo.Create(ctx, t); err != nil {
			return nil, err
		}

		results = append(results, GeneratedToken{
			ID:         t.ID,
			PlainToken: plain,
			TokenType:  t.TokenType,
			ExpiresAt:  t.ExpiresAt,
		})
	}
	return results, nil
}

func (uc *SessionUseCase) ListTokens(ctx context.Context, sessionID uuid.UUID) ([]domain.ExamToken, error) {
	return uc.tokenRepo.ListBySession(ctx, sessionID)
}

// ValidateToken checks all security rules for a token
func (uc *SessionUseCase) ValidateToken(ctx context.Context, plainToken string) (*domain.ExamToken, *domain.Session, error) {
	hash := hasher.HashToken(plainToken)

	token, err := uc.tokenRepo.FindByHash(ctx, hash)
	if err != nil {
		return nil, nil, errors.New("token tidak ditemukan")
	}

	if !token.IsActive {
		return nil, nil, errors.New("token tidak aktif")
	}

	if token.UsedCount >= token.MaxUsage {
		return nil, nil, errors.New("token sudah habis digunakan")
	}

	if token.ExpiresAt != nil && time.Now().After(*token.ExpiresAt) {
		return nil, nil, errors.New("token sudah expired")
	}

	session, err := uc.sessionRepo.FindByID(ctx, token.SessionID)
	if err != nil {
		return nil, nil, errors.New("sesi tidak ditemukan")
	}

	if session.IsLocked {
		return nil, nil, errors.New("sesi sudah dikunci")
	}

	activeCount, err := uc.sessionRepo.CountActiveParticipants(ctx, session.ID)
	if err != nil {
		return nil, nil, err
	}
	if activeCount >= session.MaxParticipants {
		return nil, nil, errors.New("sesi sudah penuh (maksimal 20 peserta)")
	}

	return token, session, nil
}
