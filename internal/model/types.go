package model

import "time"

type WorkflowRun struct {
	ID           int64
	Repo         string
	WorkflowID   int64
	WorkflowName string
	HeadBranch   string
	Event        string
	Status       string
	Conclusion   string
	RunAttempt   int
	RunStartedAt time.Time
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

type Job struct {
	ID           int64
	RunID        int64
	RunAttempt   int
	Repo         string
	Name         string
	Status       string
	Conclusion   string
	StartedAt    time.Time
	CompletedAt  time.Time
	RunnerOS     string
	RunnerName   string
	IsSelfHosted bool
	DurationSec  int
}

type OSCost struct {
	OS         string
	Minutes    float64
	Multiplier float64
	CostUSD    float64
	Percentage float64
}

type CostResult struct {
	TotalMinutes     float64
	BillableMinutes  float64
	TotalCostUSD     float64
	FreeTierUsed     float64
	ByOS             map[string]OSCost
	DataCompleteness float64
	Disclaimer       string
}

type WasteMetrics struct {
	TotalRuns       int
	FailedRuns      int
	FailRate        float64
	RerunCount      int
	RerunWasteMin   float64
	RerunWasteUSD   float64
	CancelledRuns   int
	CancelWasteMin  float64
	CancelWasteUSD  float64
	TotalWasteUSD   float64
	WastePercentage float64
}

type HotspotEntry struct {
	Rank        int
	Name        string
	GroupType   string
	Minutes     float64
	CostUSD     float64
	CostPct     float64
	RunCount    int
	FailRate    float64
	AvgDuration float64
	Trend       string
}

