package model

import "time"

type BillingSnapshot struct {
	Repo          string    `json:"repo"`
	Period        string    `json:"period"`
	ActualCostUSD float64   `json:"actual_cost_usd"`
	Source        string    `json:"source"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type ReconcileResult struct {
	Repo              string    `json:"repo"`
	Period            string    `json:"period"`
	EstimatedCostUSD  float64   `json:"estimated_cost_usd"`
	ActualCostUSD     float64   `json:"actual_cost_usd"`
	DeltaRatio        float64   `json:"delta_ratio"`
	CalibrationFactor float64   `json:"calibration_factor"`
	Confidence        string    `json:"confidence"`
	CreatedAt         time.Time `json:"created_at"`
}

type PolicyRun struct {
	Repo          string    `json:"repo"`
	PeriodStart   time.Time `json:"period_start"`
	PeriodEnd     time.Time `json:"period_end"`
	RuleID        string    `json:"rule_id"`
	Severity      string    `json:"severity"`
	Matched       bool      `json:"matched"`
	EvidenceKey   string    `json:"evidence_key"`
	EvidenceValue float64   `json:"evidence_value"`
	Expression    string    `json:"expression"`
	CreatedAt     time.Time `json:"created_at"`
}

type SuggestionRecord struct {
	Repo               string    `json:"repo"`
	PeriodStart        time.Time `json:"period_start"`
	PeriodEnd          time.Time `json:"period_end"`
	SuggestionType     string    `json:"suggestion_type"`
	Title              string    `json:"title"`
	EstimatedSavingUSD float64   `json:"estimated_saving_usd"`
	EvidenceJSON       string    `json:"evidence_json"`
	CreatedAt          time.Time `json:"created_at"`
}
