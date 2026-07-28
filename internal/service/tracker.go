package service

import (
	"context"
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/session"
	clientpkg "github.com/arihant-dev/anef-tracker/pkg/client"
	"github.com/arihant-dev/anef-tracker/pkg/config"
	"github.com/arihant-dev/anef-tracker/pkg/crawler"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"github.com/arihant-dev/anef-tracker/pkg/diff"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/event"
	"github.com/arihant-dev/anef-tracker/pkg/explorer"
	"github.com/arihant-dev/anef-tracker/pkg/log"
	"github.com/arihant-dev/anef-tracker/pkg/notify"
	"github.com/arihant-dev/anef-tracker/pkg/providers"
	v1 "github.com/arihant-dev/anef-tracker/pkg/providers/anef/v1"
	"github.com/arihant-dev/anef-tracker/pkg/recorder"
	"github.com/arihant-dev/anef-tracker/pkg/replay"
	"github.com/arihant-dev/anef-tracker/pkg/schema"
	"github.com/arihant-dev/anef-tracker/pkg/snapshot"
	"github.com/arihant-dev/anef-tracker/pkg/storage"
	"net/http"
)

type TrackerService struct {
	Config    *config.Config
	Session   *session.Session
	Store     storage.Store
	Provider  providers.Provider
	EventEng  *event.EventEngine
	SchemaEng *schema.DiscoveryEngine
	Registry  *schema.Registry
	Crawler   *crawler.Crawler
	Explorer  *explorer.Explorer
	Replayer  *replay.ReplayEngine
}

func NewTrackerService() (*TrackerService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed loading config: %w", err)
	}

	database, err := db.InitDB()
	if err != nil {
		return nil, fmt.Errorf("database init failed: %w", err)
	}

	store := storage.NewSQLiteStore(database)

	session, err := session.LoadSession()
	if err != nil {
		log.Debug("No active session file found", "error", err)
	}

	var notifiers []notify.Notifier
	if cfg.Notifications.Desktop {
		notifiers = append(notifiers, &notify.DesktopNotifier{})
	}
	if cfg.Notifications.WebhookURL != "" {
		notifiers = append(notifiers, notify.NewWebhookNotifier(cfg.Notifications.WebhookURL))
	}

	multiNotifier := notify.NewMultiNotifier(notifiers...)
	eventEngine := event.NewEventEngine(database, multiNotifier)
	schemaRegistry := schema.NewRegistry(database)
	schemaDiscovery := schema.NewDiscoveryEngine(schemaRegistry, eventEngine)

	rec := recorder.NewHTTPRecorder(database)

	authMiddleware := clientpkg.NewAuthMiddleware(http.DefaultTransport, session)
	recordingTransport := clientpkg.NewRecordingTransport(authMiddleware, rec)

	httpClient := &http.Client{Transport: recordingTransport}

	provider := v1.NewProviderV1(httpClient, rec, session)
	crl := crawler.NewCrawler(database)
	exp := explorer.NewExplorer(database)
	rep := replay.NewReplayEngine(database, httpClient)

	return &TrackerService{
		Config:    cfg,
		Session:   session,
		Store:     store,
		Provider:  provider,
		EventEng:  eventEngine,
		SchemaEng: schemaDiscovery,
		Registry:  schemaRegistry,
		Crawler:   crl,
		Explorer:  exp,
		Replayer:  rep,
	}, nil
}

// Fetch executes a full tracking cycle: fetch raw API, store snapshot, analyze schema, run diff, emit events, notify.
func (s *TrackerService) Fetch(ctx context.Context) (*domain.Application, *snapshot.SnapshotRef, *diff.DiffResult, error) {
	log.Info("Executing application status fetch via Provider", "provider", s.Provider.Name())

	app, err := s.Provider.Fetch(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("provider fetch failed: %w", err)
	}

	_ = s.Store.SaveApplication(app)

	targetURL := "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour"
	snapRef, err := snapshot.SaveHierarchicalSnapshot(app, targetURL, 200)
	if err != nil {
		log.Warn("Failed saving snapshot to disk", "error", err)
	} else {
		_ = s.Store.SaveSnapshotRef(app.ID, snapRef.SnapshotID, snapRef.Directory)
		log.Info("Hierarchical snapshot saved", "snapshot_id", snapRef.SnapshotID, "dir", snapRef.Directory)
	}

	// Schema discovery & unknown field detection
	if s.SchemaEng != nil && len(app.RawPayload) > 0 {
		discoveredEvents, err := s.SchemaEng.AnalyzeAndRegister(app.ID, targetURL, app.RawPayload)
		if err == nil && len(discoveredEvents) > 0 {
			log.Info("Discovered new schema fields", "count", len(discoveredEvents))
		}
	}

	var diffResult *diff.DiffResult
	latest, previous, _ := s.Store.GetLatestTwoSnapshots()
	if latest != nil && previous != nil {
		diffRes, err := diff.CompareSnapshots(previous.RawBytes, latest.RawBytes)
		if err == nil && diffRes.HasChanges {
			diffResult = diffRes
			_, _ = s.EventEng.ProcessSnapshotDiff(app.ID, diffRes)
			log.Info("Snapshot diff detected changes", "severity", diffRes.Severity, "change_count", len(diffRes.Changes))
		}
	}

	return app, snapRef, diffResult, nil
}

func (s *TrackerService) GetStatus() (*domain.Application, *snapshot.SnapshotRef, error) {
	latest, _, err := s.Store.GetLatestTwoSnapshots()
	if err != nil || latest == nil {
		return nil, nil, fmt.Errorf("no snapshot recorded yet")
	}
	return latest.Metadata, latest, nil
}

func (s *TrackerService) GetDiff() (*diff.DiffResult, error) {
	latest, previous, err := s.Store.GetLatestTwoSnapshots()
	if err != nil || latest == nil {
		return nil, fmt.Errorf("no snapshot recorded yet")
	}
	if previous == nil {
		return nil, fmt.Errorf("only 1 snapshot recorded. A second snapshot is required for diff comparison")
	}
	return diff.CompareSnapshots(previous.RawBytes, latest.RawBytes)
}
