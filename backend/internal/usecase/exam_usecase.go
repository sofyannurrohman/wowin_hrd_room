package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"hrd_room/backend/internal/domain"
	"hrd_room/backend/internal/repository"

	"github.com/google/uuid"
)

type ExamUseCase struct {
	participantRepo *repository.ParticipantRepository
	answerRepo      *repository.AnswerRepository
	resultRepo      *repository.ResultRepository
	questionRepo    *repository.QuestionRepository
	moduleRepo      *repository.ModuleRepository
	sessionRepo     *repository.SessionRepository
	tokenRepo       *repository.TokenRepository
	userRepo        *repository.UserRepository
}

func NewExamUseCase(
	pr *repository.ParticipantRepository,
	ar *repository.AnswerRepository,
	rr *repository.ResultRepository,
	qr *repository.QuestionRepository,
	mr *repository.ModuleRepository,
	sr *repository.SessionRepository,
	tr *repository.TokenRepository,
	ur *repository.UserRepository,
) *ExamUseCase {
	return &ExamUseCase{
		participantRepo: pr,
		answerRepo:      ar,
		resultRepo:      rr,
		questionRepo:    qr,
		moduleRepo:      mr,
		sessionRepo:     sr,
		tokenRepo:       tr,
		userRepo:        ur,
	}
}

type JoinRequest struct {
	Token string `json:"token" binding:"required"`
}

func (uc *ExamUseCase) Join(ctx context.Context, name, email string, age int, position string, token *domain.ExamToken, session *domain.Session) (*domain.SessionParticipant, error) {
	// Find or create user
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		user = &domain.User{
			ID:              uuid.New(),
			RoleID:          3, // Peserta
			Name:            name,
			Email:           email,
			PasswordHash:    "",
			Age:             &age,
			AppliedPosition: &position,
		}
		if err := uc.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	// Check if already joined
	existing, err := uc.participantRepo.FindByUserAndSession(ctx, user.ID, session.ID)
	if err == nil && existing != nil {
		// Already joined, return existing
		return existing, nil
	}

	p := &domain.SessionParticipant{
		ID:        uuid.New(),
		SessionID: session.ID,
		UserID:    user.ID,
		TokenID:   token.ID,
		JoinedAt:  time.Now(),
		Status:    "active",
	}

	if err := uc.participantRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	// Increment token usage
	if err := uc.tokenRepo.IncrementUsage(ctx, token.ID); err != nil {
		return nil, err
	}

	return p, nil
}

func (uc *ExamUseCase) GetModulesForParticipant(ctx context.Context, sessionID uuid.UUID) ([]domain.Module, error) {
	return uc.moduleRepo.ListBySession(ctx, sessionID)
}

func (uc *ExamUseCase) GetQuestionsForModule(ctx context.Context, sessionID uuid.UUID, moduleID uuid.UUID) ([]domain.Question, error) {
	session, err := uc.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	questions, err := uc.questionRepo.ListByModule(ctx, moduleID, session.RandomizeQuestions)
	if err != nil {
		return nil, err
	}

	// Mask correct answers from participants
	for i := range questions {
		for j := range questions[i].Options {
			questions[i].Options[j].IsCorrect = false
		}
	}

	return questions, nil
}

type AnswerSubmit struct {
	QuestionID       string  `json:"question_id" binding:"required"`
	SelectedOptionID *string `json:"selected_option_id"`
	TextAnswer       *string `json:"text_answer"`
}

type SubmitRequest struct {
	ParticipantID string         `json:"participant_id" binding:"required"`
	Answers       []AnswerSubmit `json:"answers" binding:"required"`
}

func (uc *ExamUseCase) SaveAnswers(ctx context.Context, participantID uuid.UUID, submissions []AnswerSubmit) error {
	for _, sub := range submissions {
		qID, err := uuid.Parse(sub.QuestionID)
		if err != nil {
			continue
		}

		a := &domain.Answer{
			ID:            uuid.New(),
			ParticipantID: participantID,
			QuestionID:    qID,
			TextAnswer:    sub.TextAnswer,
		}

		if sub.SelectedOptionID != nil {
			optID, err := uuid.Parse(*sub.SelectedOptionID)
			if err == nil {
				a.SelectedOptionID = &optID
			}
		}

		if err := uc.answerRepo.Upsert(ctx, a); err != nil {
			return err
		}
	}
	return nil
}

func (uc *ExamUseCase) AutoGrade(ctx context.Context, participantID uuid.UUID) (*domain.Result, error) {
	answers, err := uc.answerRepo.ListByParticipant(ctx, participantID)
	if err != nil {
		return nil, err
	}

	p, err := uc.participantRepo.FindByID(ctx, participantID)
	if err != nil {
		return nil, err
	}

	session, err := uc.sessionRepo.FindByID(ctx, p.SessionID)
	if err != nil {
		return nil, err
	}

	modules, err := uc.moduleRepo.ListBySession(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	var questions []domain.Question
	for _, m := range modules {
		mqs, err := uc.questionRepo.ListByModule(ctx, m.ID, false)
		if err != nil {
			return nil, err
		}
		questions = append(questions, mqs...)
	}

	questionMap := make(map[uuid.UUID]*domain.Question)
	for i := range questions {
		questionMap[questions[i].ID] = &questions[i]
	}

	var totalScore float64
	needsManualReview := false

	for _, ans := range answers {
		q, ok := questionMap[ans.QuestionID]
		if !ok {
			continue
		}

		if q.Type == domain.QuestionTypeShortAnswer && len(q.Options) > 0 {
			correct := false
			if ans.TextAnswer != nil {
				ansText := strings.TrimSpace(strings.ToLower(*ans.TextAnswer))
				for _, opt := range q.Options {
					if opt.IsCorrect {
						optText := strings.TrimSpace(strings.ToLower(opt.Content))
						if ansText == optText {
							correct = true
							break
						}
					}
				}
			}

			score := 0.0
			if correct {
				score = q.Weight
				totalScore += score
			}

			if err := uc.answerRepo.UpdateGrading(ctx, ans.ID, correct, score); err != nil {
				return nil, err
			}
			continue
		}

		if q.RequiresManualReview {
			needsManualReview = true
			continue
		}

		if ans.SelectedOptionID == nil {
			if err := uc.answerRepo.UpdateGrading(ctx, ans.ID, false, 0); err != nil {
				return nil, err
			}
			continue
		}

		// Check if selected option is correct
		correct := false
		for _, opt := range q.Options {
			if opt.ID == *ans.SelectedOptionID && opt.IsCorrect {
				correct = true
				break
			}
		}

		score := 0.0
		if correct {
			score = q.Weight
			totalScore += score
		}

		if err := uc.answerRepo.UpdateGrading(ctx, ans.ID, correct, score); err != nil {
			return nil, err
		}
	}

	gradingStatus := "completed"
	if needsManualReview {
		gradingStatus = "pending_review"
	}

	// Update participant status
	_ = uc.participantRepo.UpdateStatus(ctx, participantID, "finished")

	// Check if result already exists
	existing, err := uc.resultRepo.FindByParticipant(ctx, participantID)
	if err == nil && existing != nil {
		if err := uc.resultRepo.UpdateScore(ctx, existing.ID, totalScore, gradingStatus); err != nil {
			return nil, err
		}
		existing.TotalScore = totalScore
		existing.GradingStatus = gradingStatus
		return existing, nil
	}

	result := &domain.Result{
		ID:            uuid.New(),
		ParticipantID: participantID,
		TotalScore:    totalScore,
		GradingStatus: gradingStatus,
	}

	if err := uc.resultRepo.Create(ctx, result); err != nil {
		return nil, err
	}
	return result, nil
}

type HRReviewRequest struct {
	AnswerID string  `json:"answer_id" binding:"required"`
	Score    float64 `json:"score" binding:"required,min=0"`
	Notes    string  `json:"notes"`
}

func (uc *ExamUseCase) HRReview(ctx context.Context, req HRReviewRequest) error {
	ansID, err := uuid.Parse(req.AnswerID)
	if err != nil {
		return errors.New("invalid answer_id")
	}
	return uc.answerRepo.UpdateHRReview(ctx, ansID, req.Score, req.Notes)
}

func (uc *ExamUseCase) FinalizeScore(ctx context.Context, resultID uuid.UUID) error {
	// Fetch the result row by its own ID (not participant_id)
	result, err := uc.resultRepo.FindByID(ctx, resultID)
	if err != nil {
		return err
	}

	answers, err := uc.answerRepo.ListByParticipant(ctx, result.ParticipantID)
	if err != nil {
		return err
	}

	var total float64
	for _, a := range answers {
		if a.Score != nil {
			total += *a.Score
		}
	}

	return uc.resultRepo.UpdateScore(ctx, resultID, total, "completed")
}

func (uc *ExamUseCase) GetResults(ctx context.Context, sessionID uuid.UUID) ([]domain.Result, error) {
	return uc.resultRepo.ListBySession(ctx, sessionID)
}

func (uc *ExamUseCase) GetParticipantAnswers(ctx context.Context, participantID uuid.UUID) ([]domain.Answer, error) {
	return uc.answerRepo.ListByParticipant(ctx, participantID)
}

func (uc *ExamUseCase) GetParticipants(ctx context.Context, sessionID uuid.UUID) ([]domain.SessionParticipant, error) {
	return uc.participantRepo.ListBySession(ctx, sessionID)
}

func (uc *ExamUseCase) GetParticipantByID(ctx context.Context, participantID uuid.UUID) (*domain.SessionParticipant, error) {
	return uc.participantRepo.FindByID(ctx, participantID)
}

func (uc *ExamUseCase) UpdateSelfieUrl(ctx context.Context, participantID uuid.UUID, url string) error {
	return uc.participantRepo.UpdateSelfie(ctx, participantID, url)
}
