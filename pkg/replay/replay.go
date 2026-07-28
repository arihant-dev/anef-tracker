package replay

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type ReplayResult struct {
	ID                 int64  `json:"id"`
	OriginalRequestID  int64  `json:"original_request_id"`
	Method             string `json:"method"`
	URL                string `json:"url"`
	OriginalStatusCode int    `json:"original_status_code"`
	ReplayedStatusCode int    `json:"replayed_status_code"`
	OriginalHash       string `json:"original_hash"`
	ReplayedHash       string `json:"replayed_hash"`
	Matched            bool   `json:"matched"`
	ResponseBody       []byte `json:"response_body"`
}

type ReplayEngine struct {
	DB         *db.DB
	HTTPClient *http.Client
}

func NewReplayEngine(database *db.DB, client *http.Client) *ReplayEngine {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &ReplayEngine{
		DB:         database,
		HTTPClient: client,
	}
}

func (e *ReplayEngine) ReplayRequest(ctx context.Context, httpLogID int64, sess *session.Session) (*ReplayResult, error) {
	if e.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var method, reqURL, respBody string
	var origStatusCode int
	err := e.DB.Conn.QueryRow("SELECT method, url, status_code, resp_body FROM http_logs WHERE id = ?", httpLogID).
		Scan(&method, &reqURL, &origStatusCode, &respBody)

	if err != nil {
		return nil, fmt.Errorf("http log entry #%d not found: %w", httpLogID, err)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating replay request: %w", err)
	}

	if sess != nil {
		session.InjectAuthHeaders(req, sess)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/27.0 Safari/605.1.15")
	}

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("replay request failed: %w", err)
	}
	defer resp.Body.Close()

	replayedBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading replay response: %w", err)
	}

	origHash := HashBytes([]byte(respBody))
	replayedHash := HashBytes(replayedBytes)
	matched := (origStatusCode == resp.StatusCode) && (origHash == replayedHash)

	res := &ReplayResult{
		OriginalRequestID:  httpLogID,
		Method:             method,
		URL:                reqURL,
		OriginalStatusCode: origStatusCode,
		ReplayedStatusCode: resp.StatusCode,
		OriginalHash:       origHash,
		ReplayedHash:       replayedHash,
		Matched:            matched,
		ResponseBody:       replayedBytes,
	}

	lastID, err := e.saveReplayRecord(httpLogID, resp.StatusCode, replayedHash, matched)
	if err == nil {
		res.ID = lastID
	}

	return res, nil
}

func (e *ReplayEngine) saveReplayRecord(origID int64, statusCode int, hashStr string, matched bool) (int64, error) {
	res, err := e.DB.Conn.Exec(
		"INSERT INTO http_replays (original_request_id, timestamp, status_code, response_hash, matched) VALUES (?, ?, ?, ?, ?)",
		origID, time.Now(), statusCode, hashStr, matched,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func HashBytes(b []byte) string {
	if len(b) == 0 {
		return "EMPTY"
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])[:12]
}
