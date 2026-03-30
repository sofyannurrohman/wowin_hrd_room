package email

import (
	"fmt"
	"net/smtp"
	"time"
)

type Sender interface {
	SendInvite(toEmail, participantName, token, sessionName, loginURL string, startTime time.Time) error
}

type SMTPSender struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewSMTPSender(host, port, user, password, from string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

// Ensure SMTPSender implements Sender
var _ Sender = (*SMTPSender)(nil)

func (s *SMTPSender) SendInvite(toEmail, participantName, token, sessionName, loginURL string, startTime time.Time) error {
	if s.host == "" {
		// Log or return error if SMTP is not configured
		fmt.Printf("Warning: SMTP host is not configured. Email to %s was not sent.\n", toEmail)
		return nil // Returning nil to not break the flow if SMTP isn't set up yet
	}

	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.password, s.host)
	}

	subject := fmt.Sprintf("Undangan Ujian: %s", sessionName)

	// Convert to Asian/Jakarta (WIB) for consistency
	jakartaLoc := time.FixedZone("WIB", 7*3600)
	localTime := startTime.In(jakartaLoc)
	formattedTime := localTime.Format("02 Jan 2006 15:04") + " WIB"

	body := fmt.Sprintf(`Halo %s,

Anda telah diundang untuk mengikuti ujian: %s.

Waktu Pelaksanaan: %s

Silakan klik tautan di bawah ini untuk masuk ke halaman persiapan ujian Anda:
%s

Pada halaman tersebut masukkan Token Akses berikut ke dalam kolom "Token Ujian Khusus":
%s

(Pastikan untuk menggunakan alamat Email yang sama dengan email ini: %s)

Harap perhatikan ketentuan ujian dan persiapkan diri Anda.
Semoga berhasil!

Salam,
Tim HRD
`, participantName, sessionName, formattedTime, loginURL, token, toEmail)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", toEmail, subject, body))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, msg)
}

// MockSender is used for testing or when SMTP is disabled
type MockSender struct{}

func NewMockSender() *MockSender {
	return &MockSender{}
}

func (m *MockSender) SendInvite(toEmail, participantName, token, sessionName, loginURL string, startTime time.Time) error {
	fmt.Printf("[MOCK EMAIL] Invitation sent to %s with token %s for session %s (scheduled at: %s)\n", toEmail, token, sessionName, startTime.Format(time.RFC3339))
	return nil
}
