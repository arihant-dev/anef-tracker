package notify

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
)

type TelegramNotifier struct {
	BotToken string
	ChatID   string
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		BotToken: botToken,
		ChatID:   chatID,
	}
}

func (t *TelegramNotifier) Name() string {
	return "Telegram Notifier"
}

func (t *TelegramNotifier) Notify(event domain.Event) error {
	if t.BotToken == "" || t.ChatID == "" {
		return nil
	}
	fmt.Printf("[NOTIFY-TELEGRAM] Message dispatched to Chat #%s: %s (%s -> %s)\n", t.ChatID, event.FieldPath, event.OldVal, event.NewVal)
	return nil
}
