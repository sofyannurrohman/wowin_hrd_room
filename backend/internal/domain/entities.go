package domain

import (
	"time"

	"github.com/google/uuid"
)

// ─── Roles ──────────────────────────────────────────────────────────────────

const (
	RoleSuperAdmin = "Super Admin"
	RoleHR         = "HR"
	RolePeserta    = "Peserta"
)

// ─── User ────────────────────────────────────────────────────────────────────

type User struct {
	ID              uuid.UUID `json:"id"`
	RoleID          int       `json:"role_id"`
	RoleName        string    `json:"role_name,omitempty"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Age             *int      `json:"age,omitempty"`
	AppliedPosition *string   `json:"applied_position,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// ─── Session ─────────────────────────────────────────────────────────────────

type Session struct {
	ID                 uuid.UUID `json:"id"`
	CreatedBy          uuid.UUID `json:"created_by"`
	CreatedByName      string    `json:"created_by_name,omitempty"`
	Name               string    `json:"name"`
	Schedule           time.Time `json:"schedule"`
	DurationMinutes    int       `json:"duration_minutes"`
	MaxParticipants    int       `json:"max_participants"`
	RandomizeQuestions bool      `json:"randomize_questions"`
	ShowScore          bool      `json:"show_score"`
	IsLocked           bool      `json:"is_locked"`
	CreatedAt          time.Time `json:"created_at"`
	ParticipantCount   int       `json:"participant_count,omitempty"`
}

// ─── ExamToken ───────────────────────────────────────────────────────────────

type ExamToken struct {
	ID         uuid.UUID  `json:"id"`
	SessionID  uuid.UUID  `json:"session_id"`
	TokenHash  string     `json:"-"`
	TokenPlain string     `json:"token,omitempty"` // only returned at creation
	TokenType  string     `json:"token_type"`
	MaxUsage   int        `json:"max_usage"`
	UsedCount  int        `json:"used_count"`
	ExpiresAt  *time.Time `json:"expires_at"`
	IsActive   bool       `json:"is_active"`
	BoundEmail *string    `json:"bound_email"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ─── SessionParticipant ──────────────────────────────────────────────────────

type SessionParticipant struct {
	ID             uuid.UUID  `json:"id"`
	SessionID      uuid.UUID  `json:"session_id"`
	UserID         uuid.UUID  `json:"user_id"`
	UserName       string     `json:"user_name,omitempty"`
	UserEmail      string     `json:"user_email,omitempty"`
	TokenID        uuid.UUID  `json:"token_id"`
	JoinedAt       time.Time  `json:"joined_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
	Status         string     `json:"status"` // active, disconnected, finished
}

// ─── Module ──────────────────────────────────────────────────────────────────

type Module struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type SessionModule struct {
	SessionID uuid.UUID `json:"session_id"`
	ModuleID  uuid.UUID `json:"module_id"`
	SortOrder int       `json:"sort_order"`
}

// ─── Question ────────────────────────────────────────────────────────────────

const (
	QuestionTypeMultipleChoice = "Multiple Choice"
	QuestionTypeTrueFalse      = "True/False"
	QuestionTypeShortAnswer    = "Short Answer"
	QuestionTypePsychological  = "Psychological"
)

type Question struct {
	ID                   uuid.UUID        `json:"id"`
	ModuleID             uuid.UUID        `json:"module_id"`
	Type                 string           `json:"type"`
	Content              string           `json:"content"`
	ImageURL             *string          `json:"image_url"`
	Weight               float64          `json:"weight"`
	RequiresManualReview bool             `json:"requires_manual_review"`
	CreatedAt            time.Time        `json:"created_at"`
	Options              []QuestionOption `json:"options,omitempty"`
}

type QuestionOption struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Content    string    `json:"content"`
	IsCorrect  bool      `json:"is_correct"`
}

// ─── Answer ──────────────────────────────────────────────────────────────────

type Answer struct {
	ID               uuid.UUID  `json:"id"`
	ParticipantID    uuid.UUID  `json:"participant_id"`
	QuestionID       uuid.UUID  `json:"question_id"`
	SelectedOptionID *uuid.UUID `json:"selected_option_id"`
	TextAnswer       *string    `json:"text_answer"`
	IsCorrect        *bool      `json:"is_correct"`
	Score            *float64   `json:"score"`
	HRNotes          *string    `json:"hr_notes"`
	AnsweredAt       time.Time  `json:"answered_at"`
}

// ─── Result ──────────────────────────────────────────────────────────────────

type Result struct {
	ID            uuid.UUID `json:"id"`
	ParticipantID uuid.UUID `json:"participant_id"`
	UserName      string    `json:"user_name,omitempty"`
	UserEmail     string    `json:"user_email,omitempty"`
	TotalScore    float64   `json:"total_score"`
	GradingStatus string    `json:"grading_status"` // pending_review, completed
	SubmittedAt   time.Time `json:"submitted_at"`
}

// ─── Violation ───────────────────────────────────────────────────────────────

const (
	ViolationNoFace     = "no_face"
	ViolationMultiFace  = "multiple_face"
	ViolationLookAway   = "looking_away"
	ViolationTabSwitch  = "tab_switch"
	ViolationCameraOff  = "camera_off"
	ViolationFullscreen = "fullscreen_exit"
)

type Violation struct {
	ID            uuid.UUID `json:"id"`
	ParticipantID uuid.UUID `json:"participant_id"`
	SessionID     uuid.UUID `json:"session_id,omitempty"`
	UserName      string    `json:"user_name,omitempty"`
	ViolationType string    `json:"violation_type"`
	DetectedAt    time.Time `json:"detected_at"`
	Resolved      bool      `json:"resolved"`
}

// ─── CameraLog ───────────────────────────────────────────────────────────────

type CameraLog struct {
	ID            uuid.UUID              `json:"id"`
	ParticipantID uuid.UUID              `json:"participant_id"`
	LogData       map[string]interface{} `json:"log_data"`
	Timestamp     time.Time              `json:"timestamp"`
}

// ─── SystemLog ───────────────────────────────────────────────────────────────

type SystemLog struct {
	ID        int        `json:"id"`
	UserID    *uuid.UUID `json:"user_id"`
	UserName  string     `json:"user_name,omitempty"`
	Action    string     `json:"action"`
	Detail    string     `json:"detail"`
	CreatedAt time.Time  `json:"created_at"`
}

// ─── Job Position ────────────────────────────────────────────────────────────

type JobPosition struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
