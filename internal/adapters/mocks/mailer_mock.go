package mocks

import (
	"errors"
	"go-starter/config"
	"go-starter/internal/domain/ports"
	"sync"
)

// MailerRepositoryMock implements the ports.MailerRepository interface for testing purposes.
// It stores sent emails in memory instead of actually sending them.
type MailerRepositoryMock struct {
	cfg  *config.Mailer
	data map[string]ports.EmailMessage
	mu   sync.RWMutex
}

// NewMailerRepositoryMock creates a new instance of MailerRepositoryMock.
func NewMailerRepositoryMock(cfg *config.Mailer) *MailerRepositoryMock {
	m := &MailerRepositoryMock{
		cfg:  cfg,
		data: map[string]ports.EmailMessage{},
		mu:   sync.RWMutex{},
	}

	return m
}

// Send stores the email message in memory instead of sending it.
// The message is indexed by each recipient's email address.
func (m *MailerRepositoryMock) Send(msg ports.EmailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range msg.To {
		m.data[v] = msg
	}

	return nil
}

// Close implements the Close method of the MailerRepository interface.
// For the mock, it's a no-op operation.
func (m *MailerRepositoryMock) Close() error {
	return nil
}

// SentEmailsCount returns the total number of emails stored in the mock repository.
func (m *MailerRepositoryMock) SentEmailsCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data)
}

// GetLastSentTo retrieves the last email sent to a specific email address.
// Returns an error if no email was sent to the specified address.
func (m *MailerRepositoryMock) GetLastSentTo(email string) (ports.EmailMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	emailMessage, ok := m.data[email]
	if !ok {
		return ports.EmailMessage{}, errors.New("email not found")
	}
	return emailMessage, nil
}
