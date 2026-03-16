package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	AppEnv         string
	AppBaseURL     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBURL          string
	JWTSecret      string
	JWTExpiryHours int
	UploadDir      string
	MaxUploadMB    int64
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	maxUpload, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		AppBaseURL:     getEnv("APP_BASE_URL", "https://hrd-room.wowinapps.cloud"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "hrd_user"),
		DBPassword:     getEnv("DB_PASSWORD", "hrd_password"),
		DBName:         getEnv("DB_NAME", "hrd_room"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret"),
		JWTExpiryHours: jwtExpiry,
		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadMB:    maxUpload,
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:       getEnv("SMTP_FROM", "hrd@noreply.com"),
	}

	cfg.DBURL = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
