package mailer

import (
	"fmt"
	"go-starter/config"
	"go-starter/internal/domain/ports"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SESAdapter is an adapter for the SES service.
type SESAdapter struct {
	session   *ses.SES
	mailerCfg *config.Mailer
}

// NewSESAdapter creates a new SESAdapter instance.
func NewSESAdapter(mailerCfg *config.Mailer) (*SESAdapter, error) {
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(mailerCfg.Region),
		Credentials: credentials.NewStaticCredentials(mailerCfg.AccessKey, mailerCfg.SecretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return &SESAdapter{
		session:   ses.New(awsSession),
		mailerCfg: mailerCfg,
	}, nil
}

// Send sends an email message.
// It takes a ports.EmailMessage and returns an error if the sending fails.
func (a *SESAdapter) Send(msg ports.EmailMessage) error {
	// Convert []string to []*string for ToAddresses
	toAddresses := make([]*string, len(msg.To))
	for i, addr := range msg.To {
		toAddresses[i] = aws.String(addr)
	}

	sesInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(msg.Body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(msg.Subject),
			},
		},
		Source: aws.String(a.mailerCfg.From),
	}

	msgID, err := a.session.SendEmail(sesInput)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("email sent with message ID", "msgID", *msgID)

	return nil
}
