package services

import (
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
)

// MailerService implements the ports.MailerService interface.
// It provides high-level email sending functionality with error tracking and debug capabilities.
type MailerService struct {
	cfg        *config.Container      // Application configuration
	repo       ports.MailerRepository // Repository for sending emails
	errTracker ports.ErrorTracker     // Error tracking service
}

// NewMailerService creates a new instance of MailerService.
func NewMailerService(
	cfg *config.Container,
	repo ports.MailerRepository,
	errTracker ports.ErrorTracker,
) *MailerService {
	return &MailerService{
		cfg:        cfg,
		repo:       repo,
		errTracker: errTracker,
	}
}

// Send sends an email message through the repository.
// In non-production environments, it modifies the message for debugging purposes.
// Returns domain.ErrMailer if sending fails or if no recipients are specified.
func (m *MailerService) Send(msg *ports.EmailMessage) error {
	if len(msg.To) == 0 {
		return domain.ErrMailer
	}

	if m.cfg.Application.Env != "production" {
		m.updateForDebug(msg)
	}

	err := m.repo.Send(*msg)
	if err != nil {
		m.errTracker.CaptureException(err)
		return domain.ErrMailer
	}

	return nil
}

// updateForDebug modifies the email message for debugging purposes.
// It adds a debug prefix to the subject, appends original recipient information,
// and redirects the email to a debug address.
func (m *MailerService) updateForDebug(msg *ports.EmailMessage) {
	msg.Subject = "[DEBUG] " + msg.Subject
	msg.Body += "<br>This message was initially addressed to:"
	for _, v := range msg.To {
		msg.Body += "<br>" + v
	}

	msg.To = []string{m.cfg.Mailer.DebugTo}
}
