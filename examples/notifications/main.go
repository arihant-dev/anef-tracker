package main

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/notify"
)

func main() {
	fmt.Println("=== ANEF Tracker Example: Notifications Dispatch ===")

	emailN := notify.NewEmailNotifier("smtp.example.com", "user@example.com")
	webhookN := notify.NewWebhookNotifier("https://example.com/webhook")

	multi := notify.NewMultiNotifier(emailN, webhookN)

	ev := domain.Event{
		FieldPath: "statut",
		OldVal:    "INSTRUCTION_EN_COURS",
		NewVal:    "TITRE_A_FABRIQUER",
		Severity:  domain.SeverityHigh,
	}

	_ = multi.Notify(ev)
	fmt.Println("Dispatched notification event across multi-notifier providers.")
}
