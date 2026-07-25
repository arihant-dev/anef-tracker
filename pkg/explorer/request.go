package explorer

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"strings"
)

func FormatRequestInspection(obs *domain.EndpointObservation, reqHeaders string) string {
	var sb strings.Builder
	sb.WriteString("=== REQUEST INSPECTION ===\n")
	sb.WriteString(fmt.Sprintf("Method:  %s\n", obs.Method))
	sb.WriteString(fmt.Sprintf("URL:     %s\n", obs.URL))
	sb.WriteString(fmt.Sprintf("Headers:\n%s\n", reqHeaders))
	return sb.String()
}
