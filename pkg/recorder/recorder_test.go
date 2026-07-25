package recorder_test

import (
	"bytes"
	"compress/gzip"
	"github.com/arihant-dev/anef-tracker/pkg/recorder"
	"io"
	"net/http"
	"testing"
)

func TestRecordAndReadDecompressGzip(t *testing.T) {
	rawText := `{"numero_demande":"9929006580","statut":"TITRE_A_FABRIQUER"}`

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, _ = gzWriter.Write([]byte(rawText))
	gzWriter.Close()

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(buf.Bytes())),
	}
	resp.Header.Set("Content-Encoding", "gzip")

	rec := recorder.NewHTTPRecorder(nil)
	bodyBytes, err := rec.RecordAndRead(nil, resp, 50)
	if err != nil {
		t.Fatalf("RecordAndRead failed: %v", err)
	}

	if string(bodyBytes) != rawText {
		t.Errorf("expected decompressed body %s, got %s", rawText, string(bodyBytes))
	}
}
