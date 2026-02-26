package github

import (
	"net/http"
	"strings"
)

// NextPageURL extracts rel="next" from GitHub Link header.
func NextPageURL(h http.Header) string {
	link := h.Get("Link")
	if link == "" {
		return ""
	}
	parts := strings.Split(link, ",")
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if strings.Contains(p, `rel="next"`) {
			left := strings.Index(p, "<")
			right := strings.Index(p, ">")
			if left >= 0 && right > left {
				return p[left+1 : right]
			}
		}
	}
	return ""
}

