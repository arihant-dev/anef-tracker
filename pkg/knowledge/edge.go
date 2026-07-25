package knowledge

type EdgeType string

const (
	EdgeTypeHasStatus      EdgeType = "HAS_STATUS"
	EdgeTypeReturns        EdgeType = "RETURNS"
	EdgeTypeAttachedTo     EdgeType = "ATTACHED_TO"
	EdgeTypeTriggeredEvent EdgeType = "TRIGGERED_EVENT"
)

type Edge struct {
	From     string       `json:"from" yaml:"from"`
	To       string       `json:"to" yaml:"to"`
	Type     EdgeType     `json:"type" yaml:"type"`
	Evidence []Provenance `json:"evidence" yaml:"evidence"`
}

func (e Edge) PrimaryProvenance() Provenance {
	if len(e.Evidence) > 0 {
		return e.Evidence[0]
	}
	return Provenance{}
}
