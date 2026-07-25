package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type StatusConfig struct {
	Label         string `yaml:"label"`
	Description   string `yaml:"description"`
	SeverityLevel string `yaml:"severity_level"`
}

type StatusDict struct {
	Statuses map[string]StatusConfig `yaml:"statuses"`
}

var globalDict *StatusDict

func LoadStatusDict(path string) (*StatusDict, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dict StatusDict
	if err := yaml.Unmarshal(data, &dict); err != nil {
		return nil, err
	}
	globalDict = &dict
	return &dict, nil
}

// MapJSONToApplication parses raw ANEF JSON response bytes into a clean Application domain model.
func MapJSONToApplication(rawJSON []byte, foreignerID string) (*Application, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(rawJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse JSON into payload map: %w", err)
	}

	app := &Application{
		ForeignerID: foreignerID,
		RawPayload:  payload,
		RawJSON:     rawJSON,
		UpdatedAt:   time.Now(),
	}

	app.NumeroDemande = getStringField(payload, "numero_demande", "numeroDemande", "id")
	app.ID = app.NumeroDemande
	if app.ID == "" {
		app.ID = "ANEF-APP-UNKNOWN"
	}

	app.ForeignerID = getStringField(payload, "foreigner_id", "numero_etranger", "login")
	if app.ForeignerID == "" && foreignerID != "" {
		app.ForeignerID = foreignerID
	}

	app.LegalCategory = getStringField(payload, "categorie_juridique", "categorieJuridique", "motif_scolaire")
	app.RegulatoryReference = getStringField(payload, "reference_reglementaire", "referenceReglementaire")
	app.ProcessingSite = getStringField(payload, "site_traitement", "siteTraitement")
	app.CollectionSite = getStringField(payload, "site_retrait", "siteRetrait")
	app.Version = getInt64Field(payload, "_version", "version")

	statusCode := getStringField(payload, "statut", "code_statut", "codeStatut")
	if statusCode == "" {
		statusCode = "UNKNOWN"
	}

	app.Status = ResolveStatus(statusCode)

	// Documents parsing (Attestation + Justificatifs)
	app.Documents = parseDocuments(payload)

	return app, nil
}

func parseDocuments(payload map[string]interface{}) []Document {
	var docs []Document

	// 1. Attestation de dépôt
	if docMap, ok := payload["attestation_depot"].(map[string]interface{}); ok {
		doc := extractDocumentFromMap(docMap, "Attestation de Dépôt", "PDF")
		docs = append(docs, doc)
	}

	// 2. Justificatifs / Supporting documents array
	if rawJust, ok := payload["justificatifs"].([]interface{}); ok {
		for idx, item := range rawJust {
			if docMap, ok := item.(map[string]interface{}); ok {
				defaultName := fmt.Sprintf("Justificatif #%d", idx+1)
				doc := extractDocumentFromMap(docMap, defaultName, "DOCUMENT")
				docs = append(docs, doc)
			}
		}
	}

	return docs
}

func extractDocumentFromMap(docMap map[string]interface{}, defaultName, defaultType string) Document {
	doc := Document{
		ID:        getStringField(docMap, "id", "uuid"),
		Name:      getStringField(docMap, "name", "libelle", "nom", "type"),
		Type:      getStringField(docMap, "type", "extension", "format"),
		URL:       getStringField(docMap, "url", "chemin", "path", "href"),
		CreatedAt: time.Now(),
	}

	// Check nested fichier object
	if fichierMap, ok := docMap["fichier"].(map[string]interface{}); ok {
		if doc.Name == "" || doc.Name == defaultName {
			doc.Name = getStringField(fichierMap, "name", "nom", "filename")
		}
		if doc.URL == "" {
			doc.URL = getStringField(fichierMap, "chemin", "url", "path")
		}
	}

	if doc.Name == "" {
		doc.Name = defaultName
	}
	if doc.Type == "" {
		doc.Type = defaultType
	}
	if doc.URL == "" {
		if doc.ID != "" {
			doc.URL = fmt.Sprintf("/api/dossier/documents/%s", doc.ID)
		} else {
			doc.URL = "/api/dossier/documents/pdf"
		}
	}

	return doc
}

func ResolveStatus(code string) ApplicationStatus {
	st := ApplicationStatus{
		Code:        code,
		Label:       code,
		Description: "ANEF Workflow state: " + code,
		Severity:    SeverityMedium,
	}

	if globalDict != nil {
		if cfg, ok := globalDict.Statuses[code]; ok {
			st.Label = cfg.Label
			st.Description = cfg.Description
			st.Severity = Severity(cfg.SeverityLevel)
			return st
		}
	}

	// Default fallback mappings
	switch code {
	case "DEMANDE_SOUMISE", "DOSSIER_DEPOSE":
		st.Label = "Application Submitted"
		st.Description = "Your application has been registered successfully."
		st.Severity = SeverityLow
	case "VALIDATION_FORMELLE", "INSTRUCTION_EN_COURS":
		st.Label = "Instruction in Progress"
		st.Description = "Your application is actively under review by an officer at your prefecture."
		st.Severity = SeverityMedium
	case "TITRE_A_FABRIQUER":
		st.Label = "Residence Permit in Production"
		st.Description = "The decision has been approved and permit is in production stage."
		st.Severity = SeverityHigh
	case "TITRE_DISPONIBLE", "CONVOCATION_GENEREE":
		st.Label = "Ready for Collection"
		st.Description = "Your residence permit card has arrived at collection site."
		st.Severity = SeverityCritical
	}

	return st
}

func getStringField(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if val, ok := m[k]; ok && val != nil {
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func getInt64Field(m map[string]interface{}, keys ...string) int64 {
	for _, k := range keys {
		if val, ok := m[k]; ok && val != nil {
			switch v := val.(type) {
			case float64:
				return int64(v)
			case int64:
				return v
			case int:
				return int64(v)
			}
		}
	}
	return 0
}
