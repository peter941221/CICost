package output

import (
	"encoding/json"
	"time"
)

type JSONEnvelope struct {
	SchemaVersion string     `json:"schema_version"`
	GeneratedAt   time.Time  `json:"generated_at"`
	ToolVersion   string     `json:"tool_version"`
	Report        ReportView `json:"report"`
}

func RenderReportJSON(v ReportView, toolVersion string) (string, error) {
	payload := JSONEnvelope{
		SchemaVersion: "1.0",
		GeneratedAt:   time.Now().UTC(),
		ToolVersion:   toolVersion,
		Report:        v,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
