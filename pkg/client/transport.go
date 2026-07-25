package client

import (
	"github.com/arihant-dev/anef-tracker/pkg/recorder"

	"net/http"
	"time"
)

// RecordingTransport wraps RoundTripper to record HTTP traffic logs.
type RecordingTransport struct {
	Next     http.RoundTripper
	Recorder *recorder.HTTPRecorder
}

func NewRecordingTransport(next http.RoundTripper, rec *recorder.HTTPRecorder) *RecordingTransport {
	if next == nil {
		next = http.DefaultTransport
	}
	return &RecordingTransport{
		Next:     next,
		Recorder: rec,
	}
}

func (t *RecordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.Next.RoundTrip(req)
	latency := time.Since(start).Milliseconds()

	if t.Recorder != nil {
		if err != nil {
			t.Recorder.Record(req, nil, nil, latency, err)
		} else {
			// Response body recording handled downstream or in Provider
			t.Recorder.Record(req, resp, nil, latency, nil)
		}
	}

	return resp, err
}
