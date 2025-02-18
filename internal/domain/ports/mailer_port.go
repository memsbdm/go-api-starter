package ports

// MailerService defines the interface for email service operations.
// It provides a high-level abstraction for sending emails.
type MailerService interface {
	// Send sends an email message.
	// It takes a pointer to EmailMessage and returns an error if the sending fails.
	Send(msg *EmailMessage) error
}

// MailerRepository defines the interface for email repository operations.
// It provides low-level email sending functionality and connection management.
type MailerRepository interface {
	// Send sends an email message.
	// It takes an EmailMessage by value and returns an error if the sending fails.
	Send(msg EmailMessage) error

	// Close closes any open connections and cleans up resources.
	// Returns an error if the cleanup fails.
	Close() error
}

// EmailMessage represents an email to be sent.
// It contains the basic elements of an email message.
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
}
