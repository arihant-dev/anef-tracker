package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

// Notifier defines the decoupled notification interface.
type Notifier interface {
	Name() string
	Notify(event domain.Event) error
}

// MultiNotifier dispatches notifications to multiple providers.
type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

func (m *MultiNotifier) Name() string {
	return "MultiNotifier"
}

func (m *MultiNotifier) Notify(event domain.Event) error {
	for _, n := range m.notifiers {
		if err := n.Notify(event); err != nil {
			log.Printf("[WARN] Notification error [%s]: %v", n.Name(), err)
		}
	}
	return nil
}

// DesktopNotifier triggers local OS desktop notifications.
type DesktopNotifier struct{}

func (d *DesktopNotifier) Name() string {
	return "Desktop Notifier"
}

func (d *DesktopNotifier) Notify(event domain.Event) error {
	title := fmt.Sprintf("ANEF Status Update [%s]", event.Severity)
	msg := fmt.Sprintf("Field %s: %s → %s", event.FieldPath, event.OldVal, event.NewVal)

	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q sound name "Glass"`, msg, title)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()
	case "linux":
		cmd := exec.Command("notify-send", title, msg)
		return cmd.Run()
	}
	return nil
}

// WebhookNotifier dispatches JSON webhooks to Discord, Slack, or Telegram.
type WebhookNotifier struct {
	WebhookURL string
	HTTPClient *http.Client
}

func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		WebhookURL: url,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (w *WebhookNotifier) Name() string {
	return "Webhook Notifier"
}

func (w *WebhookNotifier) Notify(event domain.Event) error {
	if w.WebhookURL == "" {
		return nil
	}

	payload := map[string]interface{}{
		"text": fmt.Sprintf("🚨 **ANEF Status Alert** [%s]\n**Field**: `%s`\n**Old**: `%s`\n**New**: `%s`\n**Confidence**: `%.2f`",
			event.Severity, event.FieldPath, event.OldVal, event.NewVal, event.Confidence),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := w.HTTPClient.Post(w.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook responded with HTTP %d", resp.StatusCode)
	}

	return nil
}
