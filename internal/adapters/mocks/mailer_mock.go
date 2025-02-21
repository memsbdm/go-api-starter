package mocks

import (
	"errors"
	"go-starter/internal/domain/ports"
	"sync"
)

// MailerAdapterMock implements the ports.MailerAdapter interface for testing purposes.
// It stores sent emails in memory instead of actually sending them.
type MailerAdapterMock struct {
	data map[string]ports.EmailMessage
	mu   sync.RWMutex
}

// NewMailerAdapterMock creates a new instance of MailerAdapterMock.
func NewMailerAdapterMock() *MailerAdapterMock {
	m := &MailerAdapterMock{
		data: map[string]ports.EmailMessage{},
		mu:   sync.RWMutex{},
	}

	return m
}

// Send stores the email message in memory instead of sending it.
// The message is indexed by each recipient's email address.
func (m *MailerAdapterMock) Send(msg ports.EmailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range msg.To {
		m.data[v] = msg
	}

	return nil
}

// Close implements the Close method of the ports.MailerAdapter interface.
// For the mock, it's a no-op operation.
func (m *MailerAdapterMock) Close() error {
	return nil
}

// SentEmailsCount returns the total number of emails stored in the mock repository.
func (m *MailerAdapterMock) SentEmailsCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data)
}

// GetLastSentTo retrieves the last email sent to a specific email address.
// Returns an error if no email was sent to the specified address.
func (m *MailerAdapterMock) GetLastSentTo(email string) (ports.EmailMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	emailMessage, ok := m.data[email]
	if !ok {
		return ports.EmailMessage{}, errors.New("email not found")
	}
	return emailMessage, nil
}
