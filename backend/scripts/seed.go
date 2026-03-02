package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"hrd_room/backend/config"
	"hrd_room/backend/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type SeedUser struct {
	ID       uuid.UUID
	RoleID   int
	Name     string
	Email    string
	Password string
}

func hashPwd(pwd string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(b)
}

func main() {
	cfg, _ := config.Load()
	db, err := database.NewPostgresPool(cfg.DBURL)
	if err != nil {
		log.Fatal("DB error:", err)
	}
	defer db.Close()

	ctx := context.Background()

	users := []SeedUser{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), RoleID: 1, Name: "Super Admin", Email: "admin@hrdroom.com", Password: "admin123"},
		{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), RoleID: 2, Name: "HR Manager", Email: "hr@hrdroom.com", Password: "hr1234"},
		{ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"), RoleID: 3, Name: "Budi Santoso", Email: "budi@peserta.com", Password: "peserta123"},
	}

	for _, u := range users {
		hash := hashPwd(u.Password)
		_, err := db.Exec(ctx, `
			INSERT INTO users (id, role_id, name, email, password_hash, created_at)
			VALUES ($1,$2,$3,$4,$5,$6)
			ON CONFLICT (email) DO NOTHING`,
			u.ID, u.RoleID, u.Name, u.Email, hash, time.Now())
		if err != nil {
			log.Printf("Failed to seed user %s: %v", u.Email, err)
		} else {
			log.Printf("✅ User seeded: %s (%s)", u.Email, u.Password)
		}
	}

	// Seed a sample session
	sessionID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	_, err = db.Exec(ctx, `
		INSERT INTO sessions (id, created_by, name, schedule, duration_minutes, max_participants, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO NOTHING`,
		sessionID,
		uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		"Technical Assessment Q1 2026",
		time.Now().Add(1*time.Hour),
		60, 20, time.Now())
	if err != nil {
		log.Printf("Seed session error: %v", err)
	} else {
		log.Println("✅ Sample session seeded")
	}

	// Seed questions
	seedQuestions(ctx, db, sessionID)
	log.Println("🎉 Seed completed!")
}

func seedQuestions(ctx context.Context, db *pgxpool.Pool, sessionID uuid.UUID) {
	questions := []struct {
		ID      uuid.UUID
		Type    string
		Content string
		Manual  bool
		Options []struct {
			Content   string
			IsCorrect bool
		}
	}{
		{
			ID: uuid.New(), Type: "Multiple Choice",
			Content: "Apa kepanjangan dari HTTP?",
			Options: []struct {
				Content   string
				IsCorrect bool
			}{
				{"HyperText Transfer Protocol", true},
				{"High Text Transfer Protocol", false},
				{"HyperText Transmission Protocol", false},
				{"High Transfer Text Protocol", false},
			},
		},
		{
			ID: uuid.New(), Type: "True/False",
			Content: "Go adalah bahasa pemrograman yang bersifat compiled.",
			Options: []struct {
				Content   string
				IsCorrect bool
			}{
				{"True", true},
				{"False", false},
			},
		},
		{
			ID: uuid.New(), Type: "Multiple Choice",
			Content: "Database mana yang bersifat relasional?",
			Options: []struct {
				Content   string
				IsCorrect bool
			}{
				{"MongoDB", false},
				{"Redis", false},
				{"PostgreSQL", true},
				{"Cassandra", false},
			},
		},
		{
			ID: uuid.New(), Type: "Psychological", Manual: true,
			Content: "Ceritakan situasi di mana Anda menghadapi konflik dalam tim. Bagaimana Anda mengatasinya?",
		},
		{
			ID: uuid.New(), Type: "Short Answer", Manual: true,
			Content: "Jelaskan apa yang dimaksud dengan RESTful API dan sebutkan prinsip-prinsipnya.",
		},
	}

	for _, q := range questions {
		_, err := db.Exec(ctx, `
			INSERT INTO questions (id, session_id, type, content, weight, requires_manual_review, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT (id) DO NOTHING`,
			q.ID, sessionID, q.Type, q.Content, 1.0, q.Manual, time.Now())
		if err != nil {
			log.Printf("Seed question error: %v", err)
			continue
		}

		for _, opt := range q.Options {
			optID := uuid.New()
			_, err := db.Exec(ctx, `
				INSERT INTO question_options (id, question_id, content, is_correct)
				VALUES ($1,$2,$3,$4) ON CONFLICT (id) DO NOTHING`,
				optID, q.ID, opt.Content, opt.IsCorrect)
			if err != nil {
				log.Printf("Seed option error: %v", err)
			}
		}
		fmt.Printf("  ✅ Question: [%s] %s\n", q.Type, q.Content[:min(50, len(q.Content))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
