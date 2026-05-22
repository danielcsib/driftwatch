package drift_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

func TestWebhookAlerter_Send_Success(t *testing.T) {
	var received drift.Alert

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json content-type, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wa := drift.NewWebhookAlerter(ts.URL, 3*time.Second)
	alert := drift.Alert{
		Timestamp:   time.Now().UTC(),
		ServiceName: "cache",
		Level:       drift.AlertLevelWarn,
		Message:     "drift detected in service \"cache\"",
		Drifts:      []drift.DriftDetail{{Field: "image", Wanted: "v3", Actual: "v2"}},
	}

	if err := wa.Send(alert); err != nil {
		t.Fatalf("Send() unexpected error: %v", err)
	}
	if received.ServiceName != "cache" {
		t.Errorf("expected service 'cache', got %q", received.ServiceName)
	}
	if len(received.Drifts) != 1 {
		t.Errorf("expected 1 drift, got %d", len(received.Drifts))
	}
}

func TestWebhookAlerter_Send_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wa := drift.NewWebhookAlerter(ts.URL, 3*time.Second)
	alert := drift.Alert{ServiceName: "svc", Level: drift.AlertLevelError}

	if err := wa.Send(alert); err == nil {
		t.Error("expected error for non-2xx status, got nil")
	}
}

func TestNewWebhookAlerter_DefaultTimeout(t *testing.T) {
	wa := drift.NewWebhookAlerter("http://example.com", 0)
	if wa == nil {
		t.Fatal("expected non-nil WebhookAlerter")
	}
}
