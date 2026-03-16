package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"hrd_room/backend/config"
	"hrd_room/backend/internal/delivery/http/handlers"
	ws "hrd_room/backend/internal/delivery/websocket"
	"hrd_room/backend/internal/middleware"
	"hrd_room/backend/internal/repository"
	"hrd_room/backend/internal/usecase"
	"hrd_room/backend/pkg/database"
	"hrd_room/backend/pkg/email"
	jwtpkg "hrd_room/backend/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	db, err := database.NewPostgresPool(cfg.DBURL)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Ensure upload dir
	uploadDir := cfg.UploadDir
	if err := os.MkdirAll(filepath.Join(uploadDir, "questions"), 0755); err != nil {
		log.Printf("Warning: could not create upload dir: %v", err)
	}

	// Init packages
	jwtManager := jwtpkg.NewManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	answerRepo := repository.NewAnswerRepository(db)
	resultRepo := repository.NewResultRepository(db)
	moduleRepo := repository.NewModuleRepository(db)
	violationRepo := repository.NewViolationRepository(db)
	logRepo := repository.NewLogRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	jobPositionRepo := repository.NewJobPositionRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	// Use cases
	var emailSender email.Sender
	if cfg.SMTPHost != "" {
		emailSender = email.NewSMTPSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
	} else {
		emailSender = email.NewMockSender()
	}

	authUC := usecase.NewAuthUseCase(userRepo, jwtManager)
	sessionUC := usecase.NewSessionUseCase(sessionRepo, moduleRepo, tokenRepo, logRepo, emailSender, cfg.AppBaseURL)
	moduleUC := usecase.NewModuleUseCase(moduleRepo)
	jobPositionUC := usecase.NewJobPositionUseCase(jobPositionRepo)
	examUC := usecase.NewExamUseCase(participantRepo, answerRepo, resultRepo, questionRepo, moduleRepo, sessionRepo, tokenRepo, userRepo)
	settingsUC := usecase.NewSettingsUseCase(settingsRepo)

	// WebSocket hub
	hub := ws.NewHub()

	// Handlers
	authH := handlers.NewAuthHandler(authUC, jwtManager)
	sessionH := handlers.NewSessionHandler(sessionUC, hub, uploadDir)
	moduleH := handlers.NewModuleHandler(moduleUC)
	jobPositionH := handlers.NewJobPositionHandler(jobPositionUC)
	examH := handlers.NewExamHandler(examUC, sessionUC, questionRepo, violationRepo, hub, uploadDir)
	questionH := handlers.NewQuestionHTTPHandler(questionRepo, moduleRepo, uploadDir)
	adminH := handlers.NewAdminHandler(userRepo, logRepo, uploadDir)
	dashboardH := handlers.NewDashboardHandler(dashboardRepo)
	wsH := handlers.NewWSHandler(hub)
	settingsH := handlers.NewSettingsHandler(settingsUC)

	// Gin setup
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static files for uploads
	r.Static("/uploads", uploadDir)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	// ─── Public Routes ────────────────────────────────────────────────────────
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authH.Login)
			auth.POST("/register", authH.Register)
			auth.POST("/apply", authH.Apply)
		}

		// Public exam endpoints (participants have no JWT, identified by participant_id)
		api.POST("/exam/join", examH.Join)
		api.GET("/exam/:sessionId/modules", examH.GetModules)
		api.GET("/exam/:sessionId/modules/:moduleId/questions", examH.GetQuestionsForModule)
		api.POST("/exam/:sessionId/answers", examH.SubmitAnswers)
		api.POST("/exam/:sessionId/answers/autosave", examH.AutoSaveAnswers)
		api.GET("/exam/answers/:participantId", examH.GetParticipantAnswersPublic)
		api.GET("/exam/participant/:participantId/status", examH.GetParticipantStatus)
		api.POST("/violations", examH.ReportViolation)
		api.GET("/job-positions/active", jobPositionH.ListActive)
		// WebSocket (no JWT auth, uses query params for role/participant)
		api.GET("/ws/:sessionId", wsH.Handle)
	}

	// ─── Authenticated Routes ─────────────────────────────────────────────────
	authAPI := api.Group("/")
	authAPI.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Current user
		authAPI.GET("auth/me", authH.Me)
		authAPI.PUT("auth/profile", authH.UpdateProfile)
		authAPI.PUT("auth/password", authH.ChangePassword)

		// ─── HR / Admin Routes ────────────────────────────────────────────────
		hrRoutes := authAPI.Group("/")
		hrRoutes.Use(middleware.RequireRole("HR", "Super Admin"))
		{
			// Dashboard & Global Participants
			hrRoutes.GET("dashboard/stats", dashboardH.GetStats)
			hrRoutes.GET("participants", dashboardH.ListParticipants)
			hrRoutes.GET("participants/:id", dashboardH.GetParticipantDetail)

			// Sessions
			hrRoutes.GET("sessions", sessionH.List)
			hrRoutes.POST("sessions", sessionH.Create)
			hrRoutes.GET("sessions/:id", sessionH.GetByID)
			hrRoutes.PUT("sessions/:id", sessionH.Update)
			hrRoutes.PUT("sessions/:id/lock", sessionH.Lock)
			hrRoutes.DELETE("sessions/:id", sessionH.Delete)

			// Tokens
			hrRoutes.POST("sessions/:id/tokens", sessionH.GenerateTokens)
			hrRoutes.GET("sessions/:id/tokens", sessionH.ListTokens)

			// Modules
			hrRoutes.GET("modules", moduleH.List)
			hrRoutes.POST("modules", moduleH.Create)
			hrRoutes.GET("modules/:id", moduleH.GetByID)
			hrRoutes.PUT("modules/:id", moduleH.Update)
			hrRoutes.DELETE("modules/:id", moduleH.Delete)
			hrRoutes.GET("sessions/:id/modules", examH.GetModules)

			// Job Positions
			hrRoutes.GET("job-positions", jobPositionH.ListAll)
			hrRoutes.POST("job-positions", jobPositionH.Create)
			hrRoutes.GET("job-positions/:id", jobPositionH.GetByID)
			hrRoutes.PUT("job-positions/:id", jobPositionH.Update)
			hrRoutes.DELETE("job-positions/:id", jobPositionH.Delete)

			// Questions
			hrRoutes.GET("modules/:id/questions", questionH.ListByModule)
			hrRoutes.GET("questions", questionH.ListAll)
			hrRoutes.GET("questions/:id", questionH.GetByID)
			hrRoutes.POST("questions", questionH.Create)
			hrRoutes.POST("questions/import", questionH.ImportCSV)
			hrRoutes.PUT("questions/:id", questionH.Update)
			hrRoutes.DELETE("questions/:id", questionH.Delete)

			// Participants & Monitor
			hrRoutes.GET("sessions/:id/participants", examH.GetParticipants)
			hrRoutes.GET("sessions/:id/violations", examH.GetViolations)

			// Results & Review
			hrRoutes.GET("sessions/:id/results", examH.GetResults)
			hrRoutes.GET("results/:participantId/answers", examH.GetParticipantAnswers)
			hrRoutes.PUT("results/:id/review", examH.HRReview)
			hrRoutes.POST("results/:id/finalize", examH.FinalizeScore)

			// Participant user management (HR can import, update and delete participants)
			hrRoutes.POST("participants/import", adminH.ImportParticipants)
			hrRoutes.PUT("participants/:id", adminH.UpdateUser)
			hrRoutes.DELETE("participants/:id", adminH.DeleteUser)

			// Settings
			hrRoutes.GET("settings/:key", settingsH.Get)
			hrRoutes.PUT("settings/:key", settingsH.Update)
		}

		// ─── Super Admin Only ─────────────────────────────────────────────────
		adminRoutes := authAPI.Group("/admin")
		adminRoutes.Use(middleware.RequireRole("Super Admin"))
		{
			adminRoutes.GET("users", adminH.ListUsers)
			adminRoutes.PUT("users/:id", adminH.UpdateUser)
			adminRoutes.DELETE("users/:id", adminH.DeleteUser)
			adminRoutes.GET("logs", adminH.GetLogs)
		}
	}

	addr := ":" + cfg.AppPort
	log.Printf("🚀 HRD Room Backend running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(db *pgxpool.Pool) error {
	ctx := context.Background()

	// 1. Create schema_migrations table if not exists
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// 2. Read migrations folder
	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	var filenames []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			filenames = append(filenames, f.Name())
		}
	}
	sort.Strings(filenames)

	// 3. Apply missing migrations
	for _, filename := range filenames {
		var exists bool
		err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)", filename).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			log.Printf("Applying migration: %s", filename)
			content, err := os.ReadFile(filepath.Join("migrations", filename))
			if err != nil {
				return err
			}

			// Execute as one transaction per file
			tx, err := db.Begin(ctx)
			if err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, string(content)); err != nil {
				tx.Rollback(ctx)
				log.Printf("Failed to apply migration %s: %v", filename, err)
				return err
			}

			if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", filename); err != nil {
				tx.Rollback(ctx)
				return err
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
