package drift

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookAlerter sends drift alerts to an HTTP webhook endpoint.
type WebhookAlerter struct {
	url    string
	client *http.Client
}

// NewWebhookAlerter creates a WebhookAlerter that posts JSON alerts to url.
func NewWebhookAlerter(url string, timeout time.Duration) *WebhookAlerter {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &WebhookAlerter{
		url:    url,
		client: &http.Client{Timeout: timeout},
	}
}

// Send marshals the alert to JSON and POSTs it to the webhook URL.
func (w *WebhookAlerter) Send(alert Alert) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("alert_webhook: marshal: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("alert_webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("alert_webhook: unexpected status %d from %s", resp.StatusCode, w.url)
	}
	return nil
}
