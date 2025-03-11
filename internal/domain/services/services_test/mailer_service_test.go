//go:build !integration

package services_test

import (
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"reflect"
	"testing"
)

func TestMailerService_Send_Debug(t *testing.T) {
	t.Parallel()

	// Arrange
	tests := map[string]struct {
		input             *ports.EmailMessage
		expectedSent      *ports.EmailMessage
		expectedSentCount int
		expectedErr       error
	}{
		"empty to field should fail": {
			input: &ports.EmailMessage{
				To:      []string{},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSent: &ports.EmailMessage{
				To:      []string{},
				Subject: "Test",
				Body:    "Test",
			},
			expectedErr:       domain.ErrInternal,
			expectedSentCount: 0,
		},
		"send an email": {
			input: &ports.EmailMessage{
				To:      []string{"test@example.com"},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSent: &ports.EmailMessage{
				To:      []string{debugEmail},
				Subject: "[DEBUG] Test",
				Body:    "Test<br>This message was initially addressed to:<br>test@example.com",
			},
			expectedSentCount: 1,
			expectedErr:       nil,
		},
		"send an email to two addresses": {
			input: &ports.EmailMessage{
				To:      []string{"test1@example.com", "test2@example.com"},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSent: &ports.EmailMessage{
				To:      []string{debugEmail},
				Subject: "[DEBUG] Test",
				Body:    "Test<br>This message was initially addressed to:<br>test1@example.com<br>test2@example.com",
			},
			expectedSentCount: 1,
			expectedErr:       nil,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			builder := NewTestBuilder().Build()
			err := builder.MailerService.Send(tt.input)
			sentCount := getSentEmailsCount(t, builder.MailerAdapter)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if !reflect.DeepEqual(tt.input, tt.expectedSent) {
				t.Errorf("expected value %v, got %v", tt.expectedSent, tt.input)
			}
			if tt.expectedSentCount != sentCount {
				t.Errorf("expected count %v, got %v", tt.expectedSentCount, sentCount)
			}
		})
	}
}

func TestMailerService_Send_Production(t *testing.T) {
	t.Parallel()

	// Arrange
	tests := map[string]struct {
		input             *ports.EmailMessage
		expectedSent      *ports.EmailMessage
		expectedSentCount int
		expectedErr       error
	}{
		"empty to field should fail": {
			input: &ports.EmailMessage{
				To:      []string{},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSent: &ports.EmailMessage{
				To:      []string{},
				Subject: "Test",
				Body:    "Test",
			},
			expectedErr:       domain.ErrInternal,
			expectedSentCount: 0,
		},
		"send an email": {
			input: &ports.EmailMessage{
				To:      []string{"test@example.com"},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSent: &ports.EmailMessage{
				To:      []string{"test@example.com"},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSentCount: 1,
			expectedErr:       nil,
		},
		"send an email to two addresses": {
			input: &ports.EmailMessage{
				To:      []string{"test@example.com", "test2@example.com"},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSent: &ports.EmailMessage{
				To:      []string{"test@example.com", "test2@example.com"},
				Subject: "Test",
				Body:    "Test",
			},
			expectedSentCount: 2,
			expectedErr:       nil,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			builder := NewTestBuilder().SetEnvToProduction().Build()
			err := builder.MailerService.Send(tt.input)
			sentCount := getSentEmailsCount(t, builder.MailerAdapter)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if !reflect.DeepEqual(tt.input, tt.expectedSent) {
				t.Errorf("expected value %v, got %v", tt.expectedSent, tt.input)
			}
			if tt.expectedSentCount != sentCount {
				t.Errorf("expected count %v, got %v", tt.expectedSentCount, sentCount)
			}
		})
	}
}
