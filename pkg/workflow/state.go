package workflow

type State struct {
	Code        string `json:"code" yaml:"code"`
	Label       string `json:"label" yaml:"label"`
	Description string `json:"description" yaml:"description"`
}

var KnownStates = []State{
	{Code: "DEMANDE_SOUMISE", Label: "Application Submitted", Description: "Your application has been registered successfully."},
	{Code: "DOSSIER_DEPOSE", Label: "Dossier Deposited", Description: "Dossier files successfully transferred."},
	{Code: "VALIDATION_FORMELLE", Label: "Formal Validation", Description: "Formal document check passed."},
	{Code: "INSTRUCTION_EN_COURS", Label: "Instruction Started", Description: "Actively under review by an officer at prefecture."},
	{Code: "DECISION_VALIDEE", Label: "Decision Validated", Description: "Positive decision reached by prefecture."},
	{Code: "TITRE_A_FABRIQUER", Label: "Residence Permit in Production", Description: "Physical card manufacturing initiated."},
	{Code: "TITRE_DISPONIBLE", Label: "Ready for Collection", Description: "Card arrived at collection site."},
}
