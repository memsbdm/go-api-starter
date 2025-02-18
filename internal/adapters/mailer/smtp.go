package mailer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

// MailerRepository implements the ports.MailerRepository interface.
// It provides SMTP functionality with connection pooling.
type MailerRepository struct {
	cfg  *config.Mailer
	pool *sync.Pool
}

// New creates a new instance of MailerRepository with a pre-initialized pool of SMTP clients.
// It returns an error if the initialization of the SMTP client pool fails.
func New(cfg *config.Mailer) (*MailerRepository, error) {
	repo := &MailerRepository{
		cfg: cfg,
	}

	repo.pool = &sync.Pool{
		New: func() interface{} {
			client, err := createNewSMTPClient(cfg)
			if err != nil {
				slog.Error(fmt.Sprintf("error creating SMTP client: %v", err))
				return nil
			}
			return client
		},
	}

	for i := 0; i < 10; i++ {
		client, err := createNewSMTPClient(cfg)
		if err != nil {
			return nil, fmt.Errorf("error initializing pool: %w", err)
		}
		repo.pool.Put(client)
	}

	return repo, nil
}

// createNewSMTPClient creates a new SMTP client with TLS configuration.
// It handles the connection, authentication and initial testing of the SMTP connection.
func createNewSMTPClient(cfg *config.Mailer) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	tlsConfig := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("TLS connection error: %w", err)
	}
	// Defer conn close with error checking
	defer func() {
		if closeErr := conn.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing TLS connection: %w", closeErr)
		}
	}()

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("SMTP client creation error: %w", err)
	}

	defer func() {
		if err != nil {
			if closeErr := client.Close(); closeErr != nil {
				err = fmt.Errorf("multiple errors: %v; error closing SMTP client: %w", err, closeErr)
			}
		}
	}()

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err = client.Auth(auth); err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	if err := client.Noop(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	return client, nil
}

// Send attempts to send an email with retry mechanism.
// It will retry sending the email based on the configuration's MaxRetries and RetryDelayInSeconds.
func (m *MailerRepository) Send(msg ports.EmailMessage) error {
	var lastErr error
	maxRetries := m.cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	retryDelayInSeconds := m.cfg.RetryDelayInSeconds
	var retryDelay time.Duration
	if retryDelayInSeconds <= 0 {
		retryDelay = 2 * time.Second
	} else {
		retryDelay = time.Duration(retryDelayInSeconds) * time.Second
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			slog.Info(fmt.Sprintf("Attempt %d/%d sending email to %v", attempt, maxRetries, msg.To))
			time.Sleep(retryDelay)
		}

		err := m.sendWithRetry(msg)
		if err == nil {
			return nil
		}

		lastErr = err

		if !isRetryableError(err) {
			return fmt.Errorf("permanent error while sending: %w", err)
		}
	}

	return fmt.Errorf("failed after %d attempts. Last error: %w", maxRetries, lastErr)
}

// sendWithRetry handles the actual email sending process using a client from the pool.
func (m *MailerRepository) sendWithRetry(msg ports.EmailMessage) error {
	client := m.pool.Get().(*smtp.Client)
	if client == nil {
		var err error
		client, err = createNewSMTPClient(m.cfg)
		if err != nil {
			return fmt.Errorf("unable to create new client: %w", err)
		}
	}

	defer func() {
		if err := client.Noop(); err != nil {
			if closeErr := client.Close(); closeErr != nil {
				slog.Error(fmt.Sprintf("error closing SMTP client: %v", closeErr))
			}
			return
		}
		m.pool.Put(client)
	}()

	headers := make(map[string]string)
	headers["From"] = m.cfg.From
	headers["To"] = strings.Join(msg.To, ", ")
	headers["Subject"] = msg.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + msg.Body

	if err := client.Reset(); err != nil {
		return fmt.Errorf("session reset error: %w", err)
	}

	if err := client.Mail(m.cfg.From); err != nil {
		return fmt.Errorf("MAIL FROM error: %w", err)
	}

	for _, to := range msg.To {
		if err := client.Rcpt(to); err != nil {
			return fmt.Errorf("RCPT TO error: %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA error: %w", err)
	}

	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("message write error: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("writer close error: %w", err)
	}

	return nil
}

// isRetryableError determines if an error should trigger a retry attempt.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	retryableErrors := []string{
		"connection reset",
		"broken pipe",
		"connection refused",
		"no such host",
		"timeout",
		"temporary",
		"i/o timeout",
	}

	errStr := strings.ToLower(err.Error())
	for _, retryErr := range retryableErrors {
		if strings.Contains(errStr, retryErr) {
			return true
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Temporary()
	}

	return false
}

// Close properly closes all SMTP clients in the pool.
// It returns the last error encountered while closing the clients, if any.
func (m *MailerRepository) Close() error {
	var clients []*smtp.Client

	for {
		client := m.pool.Get()
		if client == nil {
			break
		}
		clients = append(clients, client.(*smtp.Client))
	}

	var lastErr error
	for _, client := range clients {
		if err := client.Quit(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
