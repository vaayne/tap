package browser

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// NetworkEntry represents a captured network request/response pair.
type NetworkEntry struct {
	RequestID    string            `json:"requestId"`
	URL          string            `json:"url"`
	Method       string            `json:"method"`
	Status       int64             `json:"status"`
	MIMEType     string            `json:"mimeType,omitempty"`
	ResourceType string            `json:"resourceType,omitempty"`
	ReqHeaders   map[string]string `json:"reqHeaders,omitempty"`
	ResHeaders   map[string]string `json:"resHeaders,omitempty"`
	Body         string            `json:"body,omitempty"`
	Error        string            `json:"error,omitempty"`
	Timestamp    time.Time         `json:"timestamp"`
}

// NetworkFilter specifies criteria for matching network requests.
type NetworkFilter struct {
	URLPattern    string   `json:"urlPattern,omitempty"`
	Methods       []string `json:"methods,omitempty"`
	ResourceTypes []string `json:"resourceTypes,omitempty"`
}

// matchURL matches a URL against a glob-like pattern where * matches any
// characters including /. An empty pattern matches everything.
func matchURL(pattern, url string) bool {
	if pattern == "" {
		return true
	}

	// Escape regex special chars, then replace escaped \* with .*
	escaped := regexp.QuoteMeta(pattern)
	rePattern := strings.ReplaceAll(escaped, `\*`, `.*`)
	rePattern = "^" + rePattern + "$"

	matched, err := regexp.MatchString(rePattern, url)
	if err != nil {
		return false
	}
	return matched
}

// matchesFilter checks whether a NetworkEntry matches the given filter.
// An empty/zero filter matches everything.
func matchesFilter(entry NetworkEntry, filter NetworkFilter) bool {
	if !matchURL(filter.URLPattern, entry.URL) {
		return false
	}

	if len(filter.Methods) > 0 {
		found := false
		for _, m := range filter.Methods {
			if strings.EqualFold(entry.Method, m) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(filter.ResourceTypes) > 0 {
		found := false
		for _, rt := range filter.ResourceTypes {
			if strings.EqualFold(entry.ResourceType, rt) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// headersToMap converts CDP Headers (map[string]any) to map[string]string.
func headersToMap(h network.Headers) map[string]string {
	if len(h) == 0 {
		return nil
	}
	out := make(map[string]string, len(h))
	for k, v := range h {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

// WaitForRequest enables the Network domain and blocks until a request matching
// the filter completes (response received + loading finished/failed). It returns
// the first matching NetworkEntry. If includeBody is true, the response body is
// fetched before returning.
//
// The caller controls the timeout via ctx (e.g. context.WithTimeout).
func WaitForRequest(ctx context.Context, debugURL string, targetID string, filter NetworkFilter, includeBody bool) (*NetworkEntry, error) {
	taskCtx, cancel, err := withTargetListen(ctx, debugURL, targetID)
	if err != nil {
		return nil, fmt.Errorf("wait for request: %w", err)
	}
	defer cancel()

	type requestState struct {
		entry    NetworkEntry
		gotResp  bool
		finished bool
		failed   bool
	}

	requests := make(map[network.RequestID]*requestState)
	done := make(chan *NetworkEntry, 1)

	chromedp.ListenTarget(taskCtx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			requests[e.RequestID] = &requestState{
				entry: NetworkEntry{
					RequestID:    string(e.RequestID),
					URL:          e.Request.URL,
					Method:       e.Request.Method,
					ResourceType: string(e.Type),
					ReqHeaders:   headersToMap(e.Request.Headers),
					Timestamp:    time.Now(),
				},
			}

		case *network.EventResponseReceived:
			rs, ok := requests[e.RequestID]
			if !ok {
				rs = &requestState{
					entry: NetworkEntry{
						RequestID: string(e.RequestID),
						Timestamp: time.Now(),
					},
				}
				requests[e.RequestID] = rs
			}
			rs.entry.Status = e.Response.Status
			rs.entry.MIMEType = e.Response.MimeType
			rs.entry.ResHeaders = headersToMap(e.Response.Headers)
			if rs.entry.URL == "" {
				rs.entry.URL = e.Response.URL
			}
			if rs.entry.ResourceType == "" {
				rs.entry.ResourceType = string(e.Type)
			}
			rs.gotResp = true

		case *network.EventLoadingFinished:
			rs, ok := requests[e.RequestID]
			if !ok {
				return
			}
			rs.finished = true
			if rs.gotResp && matchesFilter(rs.entry, filter) {
				select {
				case done <- &rs.entry:
				default:
				}
			}

		case *network.EventLoadingFailed:
			rs, ok := requests[e.RequestID]
			if !ok {
				return
			}
			rs.failed = true
			rs.entry.Error = e.ErrorText
			// Match even without ResponseReceived — a request can fail before
			// getting a response (e.g. DNS failure, connection refused). The
			// entry still has URL/Method from RequestWillBeSent and the Error
			// field captures the failure reason.
			if matchesFilter(rs.entry, filter) {
				select {
				case done <- &rs.entry:
				default:
				}
			}
		}
	})

	if err := chromedp.Run(taskCtx, network.Enable()); err != nil {
		return nil, fmt.Errorf("wait for request: enable network: %w", err)
	}

	select {
	case entry := <-done:
		if includeBody && entry.Error == "" {
			body, err := getResponseBodyOn(taskCtx, network.RequestID(entry.RequestID))
			if err == nil {
				entry.Body = string(body)
			}
			// Non-fatal: body fetch can fail for redirected/cached requests.
		}
		return entry, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for request: %w", ctx.Err())
	}
}

// GetResponseBody fetches the response body for a completed request by its ID.
func GetResponseBody(ctx context.Context, debugURL string, targetID string, requestID string) ([]byte, error) {
	var body []byte
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var fetchErr error
			body, fetchErr = getResponseBodyOn(ctx, network.RequestID(requestID))
			return fetchErr
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("get response body: %w", err)
	}
	return body, nil
}

// getResponseBodyOn fetches the response body using the given CDP execution context.
func getResponseBodyOn(ctx context.Context, requestID network.RequestID) ([]byte, error) {
	body, err := network.GetResponseBody(requestID).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("get response body %s: %w", requestID, err)
	}
	return body, nil
}
