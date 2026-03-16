package handlers

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

// handleError translates technical errors into user-friendly messages
func handleError(err error) string {
	if err == nil {
		return ""
	}

	// 1. Check for PostgreSQL specific errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			if strings.Contains(pgErr.ConstraintName, "users_email_key") || strings.Contains(pgErr.Detail, "email") {
				return "Email sudah terdaftar. Silakan gunakan email lain atau hubungi HR."
			}
			if strings.Contains(pgErr.ConstraintName, "job_positions_name_key") || strings.Contains(pgErr.Detail, "name") {
				return "Nama posisi ini sudah digunakan."
			}
			return "Data ini sudah terdaftar di sistem kami."
		case "23503": // foreign_key_violation
			return "Data tidak dapat dihapus atau diproses karena masih terhubung dengan data lain."
		}
	}

	// 2. Generic translation for common technical terms
	errMsg := err.Error()
	
	if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "23505") {
		if strings.Contains(errMsg, "email") {
			return "Email sudah terdaftar."
		}
		return "Data sudah terdaftar."
	}

	// Hide database/relation errors
	if strings.Contains(errMsg, "relation") || strings.Contains(errMsg, "sql") || strings.Contains(errMsg, "database") || strings.Contains(errMsg, "db") {
		return "Terjadi kesalahan pada sistem data. Silakan coba beberapa saat lagi atau hubungi admin."
	}

	return errMsg
}
