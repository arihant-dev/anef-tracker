package workflow

import (
	"fmt"
	"time"
)

type ObservationStrength string

const (
	StrengthHigh   ObservationStrength = "HIGH"
	StrengthMedium ObservationStrength = "MEDIUM"
	StrengthLow    ObservationStrength = "LOW"
)

type Transition struct {
	ID              int64                   `json:"id"`
	FromStatus      string                  `json:"from_status"`
	ToStatus        string                  `json:"to_status"`
	FirstSeen       time.Time               `json:"first_seen"`
	LastSeen        time.Time               `json:"last_seen"`
	AverageDuration time.Duration           `json:"average_duration"`
	Observations    []TransitionObservation `json:"observations,omitempty"`
}

func (t Transition) ObservationCount() int {
	if len(t.Observations) > 0 {
		return len(t.Observations)
	}
	return 1
}

func (t Transition) Strength() ObservationStrength {
	cnt := t.ObservationCount()
	switch {
	case cnt >= 50:
		return StrengthHigh
	case cnt >= 10:
		return StrengthMedium
	default:
		return StrengthLow
	}
}

func (t Transition) Summary() string {
	return fmt.Sprintf("%s ──> %s (Observed: %d times, Strength: %s)", t.FromStatus, t.ToStatus, t.ObservationCount(), t.Strength())
}
