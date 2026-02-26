package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(token string) *Client {
	baseURL := os.Getenv("CICOST_GITHUB_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: token,
	}
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return fmt.Sprintf("github api error (%d): %s", e.StatusCode, e.Message)
}

func (c *Client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return req, nil
}

func (c *Client) doJSON(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		msg := http.StatusText(resp.StatusCode)
		var payload struct {
			Message string `json:"message"`
		}
		if b, err2 := io.ReadAll(resp.Body); err2 == nil {
			_ = json.Unmarshal(b, &payload)
			if payload.Message != "" {
				msg = payload.Message
			}
		}
		if resp.StatusCode == http.StatusForbidden {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if sec, err3 := strconv.Atoi(retryAfter); err3 == nil && sec > 0 {
					time.Sleep(time.Duration(sec) * time.Second)
				}
			} else if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
				if epoch, err3 := strconv.ParseInt(reset, 10, 64); err3 == nil {
					sleep := time.Until(time.Unix(epoch, 0))
					if sleep > 0 && sleep < 2*time.Minute {
						time.Sleep(sleep)
					}
				}
			}
		}
		return nil, APIError{StatusCode: resp.StatusCode, Message: msg}
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}
	return resp, nil
}
