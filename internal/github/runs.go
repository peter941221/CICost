package github

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/peter941221/CICost/internal/model"
)

type runsResponse struct {
	TotalCount   int          `json:"total_count"`
	WorkflowRuns []runPayload `json:"workflow_runs"`
}

type runPayload struct {
	ID          int64  `json:"id"`
	WorkflowID  int64  `json:"workflow_id"`
	Name        string `json:"name"`
	HeadBranch  string `json:"head_branch"`
	Event       string `json:"event"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	RunAttempt  int    `json:"run_attempt"`
	RunStarted  string `json:"run_started_at"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
	WorkflowRef string `json:"path"`
}

func (c *Client) ListWorkflowRuns(ctx context.Context, owner, repo string, since time.Time) ([]model.WorkflowRun, int, error) {
	u, err := url.Parse(fmt.Sprintf("%s/repos/%s/%s/actions/runs", c.BaseURL, owner, repo))
	if err != nil {
		return nil, 0, err
	}
	q := u.Query()
	q.Set("per_page", "100")
	if !since.IsZero() {
		q.Set("created", ">="+since.UTC().Format(time.RFC3339))
	}
	u.RawQuery = q.Encode()
	nextURL := u.String()

	var out []model.WorkflowRun
	apiCalls := 0

	for nextURL != "" {
		req, err := c.newRequest(ctx, "GET", nextURL)
		if err != nil {
			return nil, apiCalls, err
		}
		var payload runsResponse
		resp, err := c.doJSON(req, &payload)
		if err != nil {
			return nil, apiCalls, err
		}
		apiCalls++
		for _, r := range payload.WorkflowRuns {
			out = append(out, mapRun(repo, r))
		}
		nextURL = NextPageURL(resp.Header)
	}
	return out, apiCalls, nil
}

func mapRun(repo string, p runPayload) model.WorkflowRun {
	return model.WorkflowRun{
		ID:           p.ID,
		Repo:         repo,
		WorkflowID:   p.WorkflowID,
		WorkflowName: p.Name,
		HeadBranch:   p.HeadBranch,
		Event:        p.Event,
		Status:       p.Status,
		Conclusion:   p.Conclusion,
		RunAttempt:   fallbackAttempt(p.RunAttempt),
		RunStartedAt: parseTime(p.RunStarted),
		UpdatedAt:    parseTime(p.UpdatedAt),
		CreatedAt:    parseTime(p.CreatedAt),
	}
}

func parseTime(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return time.Time{}
	}
	return t.UTC()
}

func fallbackAttempt(v int) int {
	if v <= 0 {
		return 1
	}
	return v
}

func parseRunID(urlRaw string) int64 {
	u, err := url.Parse(urlRaw)
	if err != nil {
		return 0
	}
	n, _ := strconv.ParseInt(u.Query().Get("run_id"), 10, 64)
	return n
}
