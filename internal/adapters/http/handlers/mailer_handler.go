package handlers

import (
	"go-starter/internal/adapters/http/responses"
	_ "go-starter/internal/adapters/http/responses"
	"go-starter/internal/domain/mailtemplates"
	"go-starter/internal/domain/ports"
	"net/http"
)

// MailerHandler is responsible for sending a test email.
type MailerHandler struct {
	errTracker ports.ErrorTracker
	mailerSvc  ports.MailerService
}

// NewMailerHandler initializes and returns a new instance of MailerHandler.
func NewMailerHandler(errTracker ports.ErrorTracker, mailerSvc ports.MailerService) *MailerHandler {
	return &MailerHandler{
		errTracker: errTracker,
		mailerSvc:  mailerSvc,
	}
}

// SendEmail godoc
//
//	@Summary		Send an example email
//	@Description	Send an example email
//	@Tags			Mail
//	@Accept			json
//	@Produce		json
//	@Success		200		"Success"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/mailer [get]
func (mh *MailerHandler) SendEmail(w http.ResponseWriter, _ *http.Request) {
	err := mh.mailerSvc.Send(&ports.EmailMessage{
		To:      []string{"example@example.com"},
		Subject: "Subject",
		Body:    mailtemplates.Hello("John Doe"),
	})
	if err != nil {
		mh.errTracker.CaptureException(err)
		responses.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
