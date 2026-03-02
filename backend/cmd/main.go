package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"hrd_room/backend/config"
	"hrd_room/backend/internal/delivery/http/handlers"
	ws "hrd_room/backend/internal/delivery/websocket"
	"hrd_room/backend/internal/middleware"
	"hrd_room/backend/internal/repository"
	"hrd_room/backend/internal/usecase"
	"hrd_room/backend/pkg/database"
	jwtpkg "hrd_room/backend/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	violationRepo := repository.NewViolationRepository(db)
	logRepo := repository.NewLogRepository(db)

	// Use cases
	authUC := usecase.NewAuthUseCase(userRepo, jwtManager)
	sessionUC := usecase.NewSessionUseCase(sessionRepo, tokenRepo, logRepo)
	examUC := usecase.NewExamUseCase(participantRepo, answerRepo, resultRepo, questionRepo, sessionRepo, tokenRepo)

	// WebSocket hub
	hub := ws.NewHub()

	// Handlers
	authH := handlers.NewAuthHandler(authUC, jwtManager)
	sessionH := handlers.NewSessionHandler(sessionUC, uploadDir)
	examH := handlers.NewExamHandler(examUC, sessionUC, questionRepo, violationRepo, uploadDir)
	questionH := handlers.NewQuestionHTTPHandler(questionRepo, uploadDir)
	adminH := handlers.NewAdminHandler(userRepo, logRepo)
	wsH := handlers.NewWSHandler(hub)

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
		}
	}

	// ─── Authenticated Routes ─────────────────────────────────────────────────
	authAPI := api.Group("/")
	authAPI.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Current user
		authAPI.GET("auth/me", authH.Me)

		// Exam (participants)
		authAPI.POST("exam/join", examH.Join)
		authAPI.GET("exam/:sessionId/questions", examH.GetQuestions)
		authAPI.POST("exam/:sessionId/answers", examH.SubmitAnswers)
		authAPI.POST("violations", examH.ReportViolation)

		// WebSocket
		authAPI.GET("ws/:sessionId", wsH.Handle)

		// ─── HR / Admin Routes ────────────────────────────────────────────────
		hrRoutes := authAPI.Group("/")
		hrRoutes.Use(middleware.RequireRole("HR", "Super Admin"))
		{
			// Sessions
			hrRoutes.GET("sessions", sessionH.List)
			hrRoutes.POST("sessions", sessionH.Create)
			hrRoutes.GET("sessions/:id", sessionH.GetByID)
			hrRoutes.PUT("sessions/:id", sessionH.Update)
			hrRoutes.DELETE("sessions/:id", sessionH.Delete)

			// Tokens
			hrRoutes.POST("sessions/:id/tokens", sessionH.GenerateTokens)
			hrRoutes.GET("sessions/:id/tokens", sessionH.ListTokens)

			// Questions
			hrRoutes.GET("sessions/:id/questions", questionH.ListBySession)
			hrRoutes.GET("questions", questionH.ListAll)
			hrRoutes.GET("questions/:id", questionH.GetByID)
			hrRoutes.POST("questions", questionH.Create)
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

func runMigrations(db interface{}) error {
	return nil // migrations handled via Docker init script
}
