package github

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/peter941221/CICost/internal/model"
)

type jobsResponse struct {
	TotalCount int          `json:"total_count"`
	Jobs       []jobPayload `json:"jobs"`
}

type jobPayload struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Conclusion      string   `json:"conclusion"`
	StartedAt       string   `json:"started_at"`
	CompletedAt     string   `json:"completed_at"`
	RunnerName      string   `json:"runner_name"`
	RunnerGroupName string   `json:"runner_group_name"`
	Labels          []string `json:"labels"`
}

func (c *Client) ListJobsForRun(ctx context.Context, owner, repo string, runID int64, attempt int) ([]model.Job, int, error) {
	apiCalls := 0
	if attempt <= 0 {
		attempt = 1
	}
	baseURL := fmt.Sprintf("%s/repos/%s/%s/actions/runs/%d/attempts/%d/jobs?per_page=100", c.BaseURL, owner, repo, runID, attempt)
	nextURL := baseURL

	var out []model.Job
	for nextURL != "" {
		req, err := c.newRequest(ctx, "GET", nextURL)
		if err != nil {
			return nil, apiCalls, err
		}
		var payload jobsResponse
		resp, err := c.doJSON(req, &payload)
		if err != nil {
			return nil, apiCalls, err
		}
		apiCalls++
		for _, j := range payload.Jobs {
			start := parseTime(j.StartedAt)
			end := parseTime(j.CompletedAt)
			dur := 0
			if !start.IsZero() && !end.IsZero() && end.After(start) {
				dur = int(end.Sub(start).Seconds())
			}
			out = append(out, model.Job{
				ID:           j.ID,
				RunID:        runID,
				RunAttempt:   attempt,
				Repo:         repo,
				Name:         j.Name,
				Status:       j.Status,
				Conclusion:   j.Conclusion,
				StartedAt:    start,
				CompletedAt:  end,
				RunnerOS:     guessRunnerOS(j.Labels, j.RunnerName),
				RunnerName:   j.RunnerName,
				RunnerGroup:  j.RunnerGroupName,
				IsSelfHosted: hasSelfHosted(j.Labels),
				DurationSec:  dur,
			})
		}
		nextURL = NextPageURL(resp.Header)
	}

	return out, apiCalls, nil
}

func hasSelfHosted(labels []string) bool {
	for _, l := range labels {
		if strings.EqualFold(l, "self-hosted") {
			return true
		}
	}
	return false
}

func guessRunnerOS(labels []string, runnerName string) string {
	for _, l := range labels {
		ll := strings.ToLower(l)
		switch {
		case strings.Contains(ll, "macos"):
			return "macOS"
		case strings.Contains(ll, "windows"):
			return "Windows"
		case strings.Contains(ll, "ubuntu"), strings.Contains(ll, "linux"):
			return "Linux"
		}
	}
	rn := strings.ToLower(runnerName)
	switch {
	case strings.Contains(rn, "mac"):
		return "macOS"
	case strings.Contains(rn, "win"):
		return "Windows"
	default:
		return "Linux"
	}
}

func extractPageFromURL(raw string) int {
	u, err := url.Parse(raw)
	if err != nil {
		return 0
	}
	return len(u.Query())
}
