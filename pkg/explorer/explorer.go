package explorer

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/crawler"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
)

type Explorer struct {
	DB      *db.DB
	Crawler *crawler.Crawler
}

func NewExplorer(database *db.DB) *Explorer {
	return &Explorer{
		DB:      database,
		Crawler: crawler.NewCrawler(database),
	}
}

func (e *Explorer) ListEndpoints() ([]domain.EndpointObservation, error) {
	return e.Crawler.DiscoverEndpoints()
}

func (e *Explorer) InspectEndpoint(id int64) (*domain.EndpointObservation, string, string, error) {
	endpoints, err := e.ListEndpoints()
	if err != nil {
		return nil, "", "", err
	}

	if id < 1 || int(id) > len(endpoints) {
		return nil, "", "", fmt.Errorf("invalid endpoint index #%d (total: %d)", id, len(endpoints))
	}

	target := endpoints[id-1]

	var reqHeaders, respBody string
	_ = e.DB.Conn.QueryRow("SELECT req_headers, resp_body FROM http_logs WHERE method = ? AND url = ? ORDER BY id DESC LIMIT 1",
		target.Method, target.URL).Scan(&reqHeaders, &respBody)

	return &target, reqHeaders, respBody, nil
}
