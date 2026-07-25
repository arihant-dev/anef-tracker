package crawler_test

import (
	"github.com/arihant-dev/anef-tracker/pkg/crawler"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"path/filepath"
	"testing"
)

func TestEndpointDiscovery(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_crawler.db")
	database, err := db.InitDBWithPath(dbPath)
	if err != nil {
		t.Fatalf("failed initializing test db: %v", err)
	}

	_ = database.RecordHTTPLog("GET", "https://example.com/api/test", 200, 150, "{}", "{}", "{}")

	crl := crawler.NewCrawler(database)
	endpoints, err := crl.DiscoverEndpoints()
	if err != nil {
		t.Fatalf("DiscoverEndpoints failed: %v", err)
	}

	if len(endpoints) != 1 {
		t.Errorf("expected 1 endpoint observation, got %d", len(endpoints))
	}

	if endpoints[0].Method != "GET" || endpoints[0].URL != "https://example.com/api/test" {
		t.Errorf("unexpected endpoint observation data")
	}
}
