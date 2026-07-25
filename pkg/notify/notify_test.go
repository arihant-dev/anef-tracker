package notify_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/notify"
	"testing"
)

func TestStatusNotification(t *testing.T) {
	emailN := notify.NewEmailNotifier("smtp.gmail.com", "user@example.com")
	if emailN.Name() != "Email Notifier" {
		t.Errorf("expected Email Notifier name")
	}

	ev := domain.Event{
		FieldPath: "statut",
		OldVal:    "INSTRUCTION_EN_COURS",
		NewVal:    "TITRE_A_FABRIQUER",
		Severity:  "HIGH",
	}

	err := emailN.Notify(ev)
	if err != nil {
		t.Errorf("expected clean notification delivery, got %v", err)
	}
}

func TestWebhookDelivery(t *testing.T) {
	webhookN := notify.NewWebhookNotifier("https://example.com/webhook")
	if webhookN.Name() != "Webhook Notifier" {
		t.Errorf("expected Webhook Notifier name")
	}
}

func TestNotificationTemplates(t *testing.T) {
	tmpl := notify.FormatStatusChangeTemplate("INSTRUCTION_EN_COURS", "TITRE_A_FABRIQUER", "snap_123", 42)
	if len(tmpl) == 0 {
		t.Errorf("expected non-empty status change template")
	}
}
