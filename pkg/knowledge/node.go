package knowledge

type NodeType string

const (
	NodeTypeApplication NodeType = "APPLICATION"
	NodeTypeStatus      NodeType = "STATUS"
	NodeTypeField       NodeType = "FIELD"
	NodeTypeDocument    NodeType = "DOCUMENT"
	NodeTypeEndpoint    NodeType = "ENDPOINT"
	NodeTypeEvent       NodeType = "EVENT"
)

type Node struct {
	ID       string                 `json:"id" yaml:"id"`
	Type     NodeType               `json:"type" yaml:"type"`
	Label    string                 `json:"label" yaml:"label"`
	Metadata map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Evidence []Provenance           `json:"evidence" yaml:"evidence"`
}

func (n Node) PrimaryProvenance() Provenance {
	if len(n.Evidence) > 0 {
		return n.Evidence[0]
	}
	return Provenance{}
}
