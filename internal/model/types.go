package model

import "time"

type WorkflowRun struct {
	ID           int64     `json:"id"`
	Repo         string    `json:"repo"`
	WorkflowID   int64     `json:"workflow_id"`
	WorkflowName string    `json:"workflow_name"`
	HeadBranch   string    `json:"head_branch"`
	Event        string    `json:"event"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	RunAttempt   int       `json:"run_attempt"`
	RunStartedAt time.Time `json:"run_started_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Job struct {
	ID           int64     `json:"id"`
	RunID        int64     `json:"run_id"`
	RunAttempt   int       `json:"run_attempt"`
	Repo         string    `json:"repo"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at"`
	RunnerOS     string    `json:"runner_os"`
	RunnerName   string    `json:"runner_name"`
	RunnerGroup  string    `json:"runner_group"`
	IsSelfHosted bool      `json:"is_self_hosted"`
	DurationSec  int       `json:"duration_sec"`
}

type OSCost struct {
	OS         string  `json:"os"`
	Minutes    float64 `json:"minutes"`
	Multiplier float64 `json:"multiplier"`
	CostUSD    float64 `json:"cost_usd"`
	Percentage float64 `json:"percentage"`
}

type CostResult struct {
	TotalMinutes     float64           `json:"total_minutes"`
	BillableMinutes  float64           `json:"billable_minutes"`
	TotalCostUSD     float64           `json:"total_cost_usd"`
	FreeTierUsed     float64           `json:"free_tier_used_min"`
	ByOS             map[string]OSCost `json:"by_os"`
	DataCompleteness float64           `json:"data_completeness"`
	Disclaimer       string            `json:"disclaimer"`
}

type WasteMetrics struct {
	TotalRuns       int     `json:"total_runs"`
	FailedRuns      int     `json:"failed_runs"`
	FailRate        float64 `json:"fail_rate"`
	RerunCount      int     `json:"rerun_count"`
	RerunWasteMin   float64 `json:"rerun_waste_min"`
	RerunWasteUSD   float64 `json:"rerun_waste_usd"`
	CancelledRuns   int     `json:"cancelled_runs"`
	CancelWasteMin  float64 `json:"cancel_waste_min"`
	CancelWasteUSD  float64 `json:"cancel_waste_usd"`
	TotalWasteUSD   float64 `json:"total_waste_usd"`
	WastePercentage float64 `json:"waste_percentage"`
}

type HotspotEntry struct {
	Rank        int     `json:"rank"`
	Name        string  `json:"name"`
	GroupType   string  `json:"group_type"`
	Minutes     float64 `json:"minutes"`
	CostUSD     float64 `json:"cost_usd"`
	CostPct     float64 `json:"cost_pct"`
	RunCount    int     `json:"run_count"`
	FailRate    float64 `json:"fail_rate"`
	AvgDuration float64 `json:"avg_duration_sec"`
	Trend       string  `json:"trend"`
}
