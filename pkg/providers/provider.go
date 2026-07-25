package providers

import (
	"context"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
)

// Provider defines the clean interface for administrative service API providers.
type Provider interface {
	Name() string
	BaseURL() string
	Fetch(ctx context.Context) (*domain.Application, error)
}
