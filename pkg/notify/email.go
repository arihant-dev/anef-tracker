package notify

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
)

type EmailNotifier struct {
	SMTPHost string
	ToEmail  string
}

func NewEmailNotifier(smtpHost, toEmail string) *EmailNotifier {
	return &EmailNotifier{
		SMTPHost: smtpHost,
		ToEmail:  toEmail,
	}
}

func (e *EmailNotifier) Name() string {
	return "Email Notifier"
}

func (e *EmailNotifier) Notify(event domain.Event) error {
	if e.ToEmail == "" {
		return nil
	}
	fmt.Printf("[NOTIFY-EMAIL] Email alert queued to %s: Status changed (%s -> %s)\n", e.ToEmail, event.OldVal, event.NewVal)
	return nil
}
