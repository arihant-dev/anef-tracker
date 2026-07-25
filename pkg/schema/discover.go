package schema

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/domain"
	"github.com/arihant-dev/anef-tracker/pkg/event"
	"log"
	"reflect"
)

type DiscoveryEngine struct {
	Registry *Registry
	EventEng *event.EventEngine
}

func NewDiscoveryEngine(reg *Registry, ee *event.EventEngine) *DiscoveryEngine {
	return &DiscoveryEngine{
		Registry: reg,
		EventEng: ee,
	}
}

// AnalyzeAndRegister inspects raw JSON payload fields, registers them in schema_fields, and emits FIELD_DISCOVERED events for new fields.
func (d *DiscoveryEngine) AnalyzeAndRegister(appID, endpoint string, payload map[string]interface{}) ([]domain.Event, error) {
	if d.Registry == nil {
		return nil, nil
	}

	var newEvents []domain.Event
	flattened := flattenMap("", payload)

	for fieldPath, val := range flattened {
		fieldType := "null"
		if val != nil {
			fieldType = reflect.TypeOf(val).String()
		}

		isNew, err := d.Registry.RegisterField(endpoint, fieldPath, fieldType)
		if err != nil {
			log.Printf("[WARN] Failed registering schema field %s: %v", fieldPath, err)
			continue
		}

		if isNew {
			valStr := fmt.Sprintf("%v", val)
			log.Printf("[SCHEMA] Discovered new field on %s: %s (%s = %s)", endpoint, fieldPath, fieldType, valStr)

			ev := domain.Event{
				ApplicationID: appID,
				Type:          "FIELD_DISCOVERED",
				Severity:      domain.SeverityLow,
				Confidence:    1.0,
				FieldPath:     fieldPath,
				OldVal:        "NONE",
				NewVal:        fmt.Sprintf("%s (%s)", valStr, fieldType),
			}

			if d.EventEng != nil && d.EventEng.DB != nil {
				_ = d.EventEng.DB.SaveEvent(ev)
			}
			newEvents = append(newEvents, ev)
		}
	}

	return newEvents, nil
}

func flattenMap(prefix string, m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		path := k
		if prefix != "" {
			path = prefix + "." + k
		}

		if subMap, ok := v.(map[string]interface{}); ok {
			subFlattened := flattenMap(path, subMap)
			for subK, subV := range subFlattened {
				result[subK] = subV
			}
		} else {
			result[path] = v
		}
	}
	return result
}
