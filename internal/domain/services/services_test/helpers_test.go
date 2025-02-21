package services_test

import (
	"go-starter/internal/domain/ports"
	"testing"
	"time"
)

func getSentEmailsCount(t *testing.T, mailer ports.MailerAdapter) int {
	t.Helper()
	if v, ok := mailer.(interface{ SentEmailsCount() int }); ok {
		return v.SentEmailsCount()
	}
	t.Fatal("the mailer adapter does not implement SentEmailsCount()")
	return 0
}

func advanceTime(t *testing.T, timeGenerator ports.TimeGenerator, duration time.Duration) {
	t.Helper()
	if v, ok := timeGenerator.(interface{ Advance(d time.Duration) }); ok {
		v.Advance(duration)
	}
}
