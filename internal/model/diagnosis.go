package model

type Severity string

const (
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type Diagnosis struct {
	RuleID         string   `json:"rule_id"`
	RuleName       string   `json:"rule_name"`
	Severity       Severity `json:"severity"`
	ResourceKind   string   `json:"resource_kind"`
	Namespace      string   `json:"namespace,omitempty"`
	ResourceName   string   `json:"resource_name"`
	ContainerName  string   `json:"container_name,omitempty"`
	Message        string   `json:"message"`
	PossibleCauses []string `json:"possible_causes,omitempty"`
	Suggestions    []string `json:"suggestions,omitempty"`
	Commands       []string `json:"commands,omitempty"`
}
