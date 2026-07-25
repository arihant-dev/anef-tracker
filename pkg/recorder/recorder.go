package recorder

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"io"
	"net/http"
	"strings"
)

type HTTPRecorder struct {
	DB *db.DB
}

func NewHTTPRecorder(database *db.DB) *HTTPRecorder {
	return &HTTPRecorder{DB: database}
}

func (r *HTTPRecorder) Record(req *http.Request, resp *http.Response, body []byte, latencyMs int64, err error) {
	if r.DB == nil {
		return
	}

	reqHeadersStr := ""
	if req != nil {
		reqHeadersStr = headerToJSON(SanitizeHeaders(req.Header))
	}

	respHeadersStr := ""
	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
		respHeadersStr = headerToJSON(SanitizeHeaders(resp.Header))
	}

	bodyStr := ""
	if len(body) > 0 {
		bodyStr = string(body)
	}

	reqURL := ""
	reqMethod := "GET"
	if req != nil {
		reqURL = req.URL.String()
		reqMethod = req.Method
	}

	_ = r.DB.RecordHTTPLog(reqMethod, reqURL, statusCode, latencyMs, reqHeadersStr, respHeadersStr, bodyStr)
}

func (r *HTTPRecorder) RecordAndRead(req *http.Request, resp *http.Response, latencyMs int64) ([]byte, error) {
	var bodyBytes []byte
	if resp != nil && resp.Body != nil {
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Decompress GZIP or Deflate if compressed
		encoding := strings.ToLower(resp.Header.Get("Content-Encoding"))
		if encoding == "gzip" || (len(raw) >= 2 && raw[0] == 0x1f && raw[1] == 0x8b) {
			gzReader, err := gzip.NewReader(bytes.NewReader(raw))
			if err == nil {
				decompressed, err := io.ReadAll(gzReader)
				gzReader.Close()
				if err == nil {
					raw = decompressed
				}
			}
		} else if encoding == "deflate" {
			flateReader := flate.NewReader(bytes.NewReader(raw))
			decompressed, err := io.ReadAll(flateReader)
			flateReader.Close()
			if err == nil {
				raw = decompressed
			}
		}

		bodyBytes = raw
	}

	r.Record(req, resp, bodyBytes, latencyMs, nil)
	return bodyBytes, nil
}

func headerToJSON(h http.Header) string {
	b, _ := json.Marshal(h)
	return string(b)
}
