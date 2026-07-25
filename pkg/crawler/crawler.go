package crawler

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"time"
)

type Crawler struct {
	DB *db.DB
}

func NewCrawler(database *db.DB) *Crawler {
	return &Crawler{DB: database}
}

// DiscoverEndpoints queries http_logs table and compiles observed endpoint statistics.
func (c *Crawler) DiscoverEndpoints() ([]domain.EndpointObservation, error) {
	if c.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT 
		method, 
		url, 
		MAX(status_code) as last_status_code, 
		MAX(latency_ms) as last_latency_ms, 
		MIN(created_at) as first_seen, 
		MAX(created_at) as last_seen, 
		COUNT(*) as occurrences 
	FROM http_logs 
	GROUP BY method, url 
	ORDER BY occurrences DESC`

	rows, err := c.DB.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []domain.EndpointObservation
	idx := int64(1)
	for rows.Next() {
		var obs domain.EndpointObservation
		obs.ID = idx
		idx++
		var firstSeenVal, lastSeenVal interface{}

		if err := rows.Scan(&obs.Method, &obs.URL, &obs.LastStatusCode, &obs.LastLatencyMs, &firstSeenVal, &lastSeenVal, &obs.Occurrences); err != nil {
			return nil, err
		}

		obs.FirstSeen = parseTime(firstSeenVal)
		obs.LastSeen = parseTime(lastSeenVal)

		observations = append(observations, obs)
	}

	return observations, nil
}

func parseTime(val interface{}) time.Time {
	if val == nil {
		return time.Now()
	}
	switch v := val.(type) {
	case time.Time:
		return v
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err == nil {
			return t
		}
		t2, err := time.Parse(time.RFC3339, v)
		if err == nil {
			return t2
		}
	}
	return time.Now()
}
