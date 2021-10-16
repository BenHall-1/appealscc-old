package mail

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

func SendEmail(email, subject, body string) string {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(mailgun.APIBaseEU)

	if err != nil {
		sentry.CaptureException(err)
	}

	m := mg.NewMessage(
		"Appeals.CC <noreply@mail.appeals.cc>",
		subject,
		body,
		email,
	)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	_, id, err := mg.Send(ctx, m)

	if err != nil {
		sentry.CaptureException(err)
	}

	return id
}
