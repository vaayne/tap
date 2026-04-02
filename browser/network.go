package browser

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/fetch"
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

// compileURLPattern compiles a glob-like URL pattern into a regexp.
// An empty pattern returns nil (matches everything).
func compileURLPattern(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil
	}
	escaped := regexp.QuoteMeta(pattern)
	rePattern := "^" + strings.ReplaceAll(escaped, `\*`, `.*`) + "$"
	return regexp.MustCompile(rePattern)
}

// matchURL matches a URL against a glob-like pattern where * matches any
// characters including /. An empty pattern matches everything.
func matchURL(pattern, url string) bool {
	re := compileURLPattern(pattern)
	if re == nil {
		return true
	}
	return re.MatchString(url)
}

// compiledFilter is a pre-compiled version of NetworkFilter for efficient
// repeated matching against many entries.
type compiledFilter struct {
	urlRe         *regexp.Regexp
	methods       []string
	resourceTypes []string
}

// compileFilter pre-compiles a NetworkFilter for repeated use.
func compileFilter(f NetworkFilter) compiledFilter {
	return compiledFilter{
		urlRe:         compileURLPattern(f.URLPattern),
		methods:       f.Methods,
		resourceTypes: f.ResourceTypes,
	}
}

// containsFold reports whether any element in slice equals s (case-insensitive).
func containsFold(slice []string, s string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}

// matches checks whether a NetworkEntry matches this compiled filter.
func (cf compiledFilter) matches(entry NetworkEntry) bool {
	if cf.urlRe != nil && !cf.urlRe.MatchString(entry.URL) {
		return false
	}
	if len(cf.methods) > 0 && !containsFold(cf.methods, entry.Method) {
		return false
	}
	if len(cf.resourceTypes) > 0 && !containsFold(cf.resourceTypes, entry.ResourceType) {
		return false
	}
	return true
}

// matchesFilter checks whether a NetworkEntry matches the given filter.
// An empty/zero filter matches everything.
func matchesFilter(entry NetworkEntry, filter NetworkFilter) bool {
	return compileFilter(filter).matches(entry)
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

// ---------------------------------------------------------------------------
// Shared request tracker
// ---------------------------------------------------------------------------

// requestState tracks the lifecycle of a single network request.
type requestState struct {
	entry   NetworkEntry
	gotResp bool
}

// requestTracker accumulates Network domain events and emits completed entries
// that match the filter via the onMatch callback.
type requestTracker struct {
	filter   compiledFilter
	onMatch  func(NetworkEntry)
	requests map[network.RequestID]*requestState
}

// newRequestTracker creates a tracker with a pre-compiled filter.
func newRequestTracker(filter NetworkFilter, onMatch func(NetworkEntry)) *requestTracker {
	return &requestTracker{
		filter:   compileFilter(filter),
		onMatch:  onMatch,
		requests: make(map[network.RequestID]*requestState),
	}
}

// handleEvent processes a single CDP Network event.
func (rt *requestTracker) handleEvent(ev interface{}) {
	switch e := ev.(type) {
	case *network.EventRequestWillBeSent:
		rt.requests[e.RequestID] = &requestState{
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
		rs, ok := rt.requests[e.RequestID]
		if !ok {
			rs = &requestState{
				entry: NetworkEntry{
					RequestID: string(e.RequestID),
					Timestamp: time.Now(),
				},
			}
			rt.requests[e.RequestID] = rs
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
		rs, ok := rt.requests[e.RequestID]
		if !ok {
			return
		}
		entry := rs.entry
		gotResp := rs.gotResp
		delete(rt.requests, e.RequestID)
		if gotResp && rt.filter.matches(entry) {
			rt.onMatch(entry)
		}

	case *network.EventLoadingFailed:
		rs, ok := rt.requests[e.RequestID]
		if !ok {
			return
		}
		entry := rs.entry
		entry.Error = e.ErrorText
		delete(rt.requests, e.RequestID)
		// Match even without ResponseReceived — a request can fail before
		// getting a response (e.g. DNS failure, connection refused).
		if rt.filter.matches(entry) {
			rt.onMatch(entry)
		}
	}
}

// ---------------------------------------------------------------------------
// Network domain primitives
// ---------------------------------------------------------------------------

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

	done := make(chan NetworkEntry, 1)
	tracker := newRequestTracker(filter, func(entry NetworkEntry) {
		select {
		case done <- entry:
		default:
		}
	})

	chromedp.ListenTarget(taskCtx, tracker.handleEvent)

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
		return &entry, nil
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

// EnableNetworkLog enables the Network domain on a target and streams completed
// request/response entries to the returned channel. The channel is buffered
// (256 entries); if the buffer fills, new entries are dropped.
//
// Call the returned cancel func to stop capturing and close the channel.
// The caller should drain the channel or call cancel to avoid goroutine leaks.
func EnableNetworkLog(ctx context.Context, debugURL string, targetID string, filter NetworkFilter) (<-chan NetworkEntry, func(), error) {
	taskCtx, taskCancel, err := withTargetListen(ctx, debugURL, targetID)
	if err != nil {
		return nil, nil, fmt.Errorf("enable network log: %w", err)
	}

	ch := make(chan NetworkEntry, 256)
	tracker := newRequestTracker(filter, func(entry NetworkEntry) {
		select {
		case ch <- entry:
		default:
			// Buffer full — drop entry.
		}
	})

	chromedp.ListenTarget(taskCtx, tracker.handleEvent)

	if err := chromedp.Run(taskCtx, network.Enable()); err != nil {
		taskCancel()
		return nil, nil, fmt.Errorf("enable network log: enable network: %w", err)
	}

	// Goroutine closes the channel when the CDP session ends.
	go func() {
		<-taskCtx.Done()
		close(ch)
	}()

	return ch, taskCancel, nil
}

// InterceptRule defines how to modify or block requests matching a filter.
// Block and MockBody are mutually exclusive.
type InterceptRule struct {
	Filter      NetworkFilter     `json:"filter"`
	Block       bool              `json:"block,omitempty"`
	AddHeaders  map[string]string `json:"addHeaders,omitempty"`
	MockStatus  int               `json:"mockStatus,omitempty"`
	MockBody    string            `json:"mockBody,omitempty"`
	MockHeaders map[string]string `json:"mockHeaders,omitempty"`
}

// ValidateInterceptRules checks that rules are well-formed.
func ValidateInterceptRules(rules []InterceptRule) error {
	for i, r := range rules {
		if r.Block && r.MockBody != "" {
			return fmt.Errorf("rule %d: block and mock body are mutually exclusive", i)
		}
		if r.MockBody != "" && r.MockStatus == 0 {
			return fmt.Errorf("rule %d: mock body requires mock status", i)
		}
	}
	return nil
}

// SetInterceptRules enables Fetch domain interception with the given rules.
// It replaces any previously set rules. The returned cancel func stops the
// interception goroutine and disables the Fetch domain.
//
// Pass nil/empty rules to effectively disable interception (or use ClearIntercept).
func SetInterceptRules(ctx context.Context, debugURL string, targetID string, rules []InterceptRule) (func(), error) {
	if err := ValidateInterceptRules(rules); err != nil {
		return nil, fmt.Errorf("set intercept rules: %w", err)
	}

	taskCtx, taskCancel, err := withTargetListen(ctx, debugURL, targetID)
	if err != nil {
		return nil, fmt.Errorf("set intercept rules: %w", err)
	}

	// Build CDP RequestPattern entries from rules.
	var patterns []*fetch.RequestPattern
	for _, r := range rules {
		p := &fetch.RequestPattern{}
		if r.Filter.URLPattern != "" {
			p.URLPattern = r.Filter.URLPattern
		}
		if len(r.Filter.ResourceTypes) == 1 {
			p.ResourceType = network.ResourceType(r.Filter.ResourceTypes[0])
		}
		patterns = append(patterns, p)
	}
	if len(patterns) == 0 {
		// Match all requests if no patterns specified.
		patterns = []*fetch.RequestPattern{{}}
	}

	// Register the event handler before enabling the domain.
	chromedp.ListenTarget(taskCtx, func(ev interface{}) {
		e, ok := ev.(*fetch.EventRequestPaused)
		if !ok {
			return
		}

		// Find the first matching rule.
		entry := NetworkEntry{
			URL:          e.Request.URL,
			Method:       e.Request.Method,
			ResourceType: string(e.ResourceType),
		}

		var matched *InterceptRule
		for i := range rules {
			if matchesFilter(entry, rules[i].Filter) {
				matched = &rules[i]
				break
			}
		}

		// Must respond to every paused request. If no rule matches, continue.
		// Run CDP commands in a goroutine to avoid blocking the event loop.
		go func() {
			if matched == nil {
				_ = fetch.ContinueRequest(e.RequestID).Do(taskCtx)
				return
			}

			if matched.Block {
				_ = fetch.FailRequest(e.RequestID, network.ErrorReasonBlockedByClient).Do(taskCtx)
				return
			}

			if matched.MockBody != "" {
				headers := []*fetch.HeaderEntry{}
				for k, v := range matched.MockHeaders {
					headers = append(headers, &fetch.HeaderEntry{Name: k, Value: v})
				}
				if len(headers) == 0 {
					headers = append(headers, &fetch.HeaderEntry{
						Name: "Content-Type", Value: "application/json",
					})
				}
				_ = fetch.FulfillRequest(e.RequestID, int64(matched.MockStatus)).
					WithResponseHeaders(headers).
					WithBody(base64.StdEncoding.EncodeToString([]byte(matched.MockBody))).
					Do(taskCtx)
				return
			}

			// AddHeaders only — continue with modified headers.
			if len(matched.AddHeaders) > 0 {
				// Build a case-insensitive map to deduplicate headers.
				// AddHeaders overrides existing headers with the same name.
				merged := make(map[string]string)
				mergedKeys := make(map[string]string) // lowercase → original case
				for k, v := range e.Request.Headers {
					lk := strings.ToLower(k)
					merged[lk] = fmt.Sprintf("%v", v)
					mergedKeys[lk] = k
				}
				for k, v := range matched.AddHeaders {
					lk := strings.ToLower(k)
					merged[lk] = v
					mergedKeys[lk] = k // prefer the casing from AddHeaders
				}
				headers := make([]*fetch.HeaderEntry, 0, len(merged))
				for lk, v := range merged {
					headers = append(headers, &fetch.HeaderEntry{Name: mergedKeys[lk], Value: v})
				}
				_ = fetch.ContinueRequest(e.RequestID).WithHeaders(headers).Do(taskCtx)
				return
			}

			// Rule matched but no action specified — continue normally.
			_ = fetch.ContinueRequest(e.RequestID).Do(taskCtx)
		}()
	})

	if err := chromedp.Run(taskCtx, fetch.Enable().WithPatterns(patterns)); err != nil {
		taskCancel()
		return nil, fmt.Errorf("set intercept rules: enable fetch: %w", err)
	}

	return taskCancel, nil
}

// ClearIntercept disables the Fetch domain on a target.
func ClearIntercept(ctx context.Context, debugURL string, targetID string) error {
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			return fetch.Disable().Do(ctx)
		}),
	)
	if err != nil {
		return fmt.Errorf("clear intercept: %w", err)
	}
	return nil
}
