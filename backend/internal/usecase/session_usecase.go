package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"
	"hrd_room/backend/pkg/email"
	"hrd_room/backend/pkg/hasher"

	"github.com/google/uuid"
)

type SessionUseCase struct {
	sessionRepo    *repository.SessionRepository
	moduleRepo     *repository.ModuleRepository
	tokenRepo      *repository.TokenRepository
	logRepo        *repository.LogRepository
	monitoringRepo *repository.MonitoringRepository
	violationRepo  *repository.ViolationRepository
	emailSender    email.Sender
	baseURL        string
	uploadDir      string
}

func NewSessionUseCase(sr *repository.SessionRepository, mr *repository.ModuleRepository, tr *repository.TokenRepository, lr *repository.LogRepository, monr *repository.MonitoringRepository, vr *repository.ViolationRepository, es email.Sender, baseURL, uploadDir string) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo:    sr,
		moduleRepo:     mr,
		tokenRepo:      tr,
		logRepo:        lr,
		monitoringRepo: monr,
		violationRepo:  vr,
		emailSender:    es,
		baseURL:        baseURL,
		uploadDir:      uploadDir,
	}
}

type CreateSessionRequest struct {
	Name               string      `json:"name" binding:"required"`
	Schedule           time.Time   `json:"schedule" binding:"required"`
	DurationMinutes    int         `json:"duration_minutes" binding:"required,min=1"`
	MaxParticipants    int         `json:"max_participants"`
	RandomizeQuestions bool        `json:"randomize_questions"`
	ShowScore          bool        `json:"show_score"`
	ModuleIDs          []uuid.UUID `json:"module_ids" binding:"required,gt=0"`
}

func (uc *SessionUseCase) Create(ctx context.Context, createdBy uuid.UUID, req CreateSessionRequest) (*domain.Session, error) {
	if len(req.ModuleIDs) == 0 {
		return nil, errors.New("minimal satu modul soal harus dipilih")
	}

	maxP := req.MaxParticipants
	if maxP <= 0 || maxP > 20 {
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

	for i, mID := range req.ModuleIDs {
		sm := &domain.SessionModule{
			SessionID: s.ID,
			ModuleID:  mID,
			SortOrder: i,
		}
		if err := uc.moduleRepo.AddSessionModule(ctx, sm); err != nil {
			// If we fail here, we should ideally rollback, but at least we don't return success
			return nil, fmt.Errorf("gagal menambahkan modul ke sesi: %w", err)
		}
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
	if len(req.ModuleIDs) == 0 {
		return nil, errors.New("minimal satu modul soal harus dipilih")
	}

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

	if err := uc.sessionRepo.Update(ctx, s); err != nil {
		return nil, err
	}

	if err := uc.moduleRepo.DeleteSessionModules(ctx, s.ID); err != nil {
		return nil, fmt.Errorf("gagal memperbarui modul sesi: %w", err)
	}

	for i, mID := range req.ModuleIDs {
		sm := &domain.SessionModule{
			SessionID: s.ID,
			ModuleID:  mID,
			SortOrder: i,
		}
		if err := uc.moduleRepo.AddSessionModule(ctx, sm); err != nil {
			return nil, fmt.Errorf("gagal menambahkan modul ke sesi: %w", err)
		}
	}

	return s, nil
}

func (uc *SessionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.sessionRepo.Delete(ctx, id)
}

func (uc *SessionUseCase) SaveMonitoringPhoto(ctx context.Context, sessionID, participantID uuid.UUID, base64Data string) error {
	path, err := uc.saveBase64Image(base64Data, "monitoring")
	if err != nil {
		return err
	}

	p := &domain.MonitoringPhoto{
		ID:            uuid.New(),
		SessionID:     sessionID,
		ParticipantID: participantID,
		PhotoURL:      path,
		CreatedAt:     time.Now(),
	}

	return uc.monitoringRepo.Create(ctx, p)
}

func (uc *SessionUseCase) ReportViolationWithProof(ctx context.Context, sessionID, participantID uuid.UUID, vType, base64Data string) error {
	var proofURL *string
	if base64Data != "" {
		path, err := uc.saveBase64Image(base64Data, "violations")
		if err == nil {
			proofURL = &path
		}
	}

	v := &domain.Violation{
		ID:            uuid.New(),
		ParticipantID: participantID,
		SessionID:     sessionID,
		ViolationType: vType,
		DetectedAt:    time.Now(),
		ProofURL:      proofURL,
	}

	return uc.violationRepo.Create(ctx, v)
}

func (uc *SessionUseCase) GetMonitoringPhotos(ctx context.Context, participantID uuid.UUID) ([]domain.MonitoringPhoto, error) {
	return uc.monitoringRepo.ListByParticipant(ctx, participantID)
}

func (uc *SessionUseCase) CleanupOldMonitoringPhotos(ctx context.Context) (int, error) {
	const days = 7
	photos, err := uc.monitoringRepo.GetOldPhotos(ctx, days)
	if err != nil {
		return 0, err
	}

	var photoIDs []uuid.UUID
	count := 0
	for _, p := range photos {
		// Try delete file
		filePath := filepath.Join(uc.uploadDir, strings.TrimPrefix(p.PhotoURL, "/uploads/"))
		_ = os.Remove(filePath)
		photoIDs = append(photoIDs, p.ID)
		count++
	}

	if len(photoIDs) > 0 {
		_ = uc.monitoringRepo.DeletePhotos(ctx, photoIDs)
	}

	// Also cleanup violation proofs
	viols, err := uc.monitoringRepo.GetOldViolations(ctx, days)
	if err == nil {
		var violIDs []uuid.UUID
		for _, v := range viols {
			if v.ProofURL != nil {
				filePath := filepath.Join(uc.uploadDir, strings.TrimPrefix(*v.ProofURL, "/uploads/"))
				_ = os.Remove(filePath)
				violIDs = append(violIDs, v.ID)
				count++
			}
		}
		if len(violIDs) > 0 {
			_ = uc.monitoringRepo.ClearProofURLs(ctx, violIDs)
		}
	}

	return count, nil
}

func (uc *SessionUseCase) saveBase64Image(data, subDir string) (string, error) {
	// data is "data:image/jpeg;base64,..."
	parts := strings.Split(data, ",")
	if len(parts) < 2 {
		return "", errors.New("invalid base64 data")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	dir := filepath.Join(uc.uploadDir, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s.jpg", time.Now().UnixNano(), uuid.New().String()[:8])
	dst := filepath.Join(dir, filename)

	if err := os.WriteFile(dst, decoded, 0644); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s/%s", subDir, filename), nil
}

func (uc *SessionUseCase) Lock(ctx context.Context, id uuid.UUID, locked bool) error {
	return uc.sessionRepo.Lock(ctx, id, locked)
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
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, errors.New("sesi tidak ditemukan")
	}

	if session.IsLocked {
		return nil, errors.New("tidak dapat generate token: sesi ujian sudah dikunci")
	}

	if req.BoundEmail == nil || *req.BoundEmail == "" {
		return nil, errors.New("email wajib diisi untuk keamanan (token harus tertaut ke peserta)")
	}

	endTime := session.Schedule.Add(time.Duration(session.DurationMinutes) * time.Minute)
	if time.Now().After(endTime) {
		return nil, errors.New("tidak dapat generate token: sesi ujian sudah selesai")
	}

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

		// Email logic implemented here:
		if t.BoundEmail != nil && *t.BoundEmail != "" {
			session, _ := uc.sessionRepo.FindByID(ctx, sessionID) // We can err check silently
			sessionName := "Ujian HRD"
			startTime := time.Now()
			if session != nil {
				sessionName = session.Name
				startTime = session.Schedule
			}

			// Try to find the participant name. In our context this isn't strictly necessary for bare functionality,
			// but could be enhanced. Using empty string or email for now until deeper user lookups.

			// Build magic login link — reads APP_BASE_URL env var (defaults to production domain)
			loginURL := uc.baseURL + "/join"

			go uc.emailSender.SendInvite(*t.BoundEmail, "Peserta", plain, sessionName, loginURL, startTime)
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

	if time.Now().Before(session.Schedule) {
		return nil, nil, errors.New("sesi ujian belum dimulai")
	}

	endTime := session.Schedule.Add(time.Duration(session.DurationMinutes) * time.Minute)
	if time.Now().After(endTime) {
		return nil, nil, errors.New("sesi ujian sudah selesai")
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
