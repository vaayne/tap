package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fastschema/qjs"
	"github.com/vaayne/tap/script"
	"github.com/vaayne/tap/transport"
)

// QuickJS executes scripts in an embedded QuickJS runtime with a Go-backed fetch().
type QuickJS struct {
	transport *transport.Transport
}

// NewQuickJS creates a new QuickJS engine backed by the given transport.
func NewQuickJS(tp *transport.Transport) *QuickJS {
	return &QuickJS{transport: tp}
}

func (q *QuickJS) Name() string { return "QuickJS" }
func (q *QuickJS) Close() error { return nil }

func (q *QuickJS) Run(_ context.Context, s *script.Script, args map[string]string) (any, error) {
	rt, err := qjs.New()
	if err != nil {
		return nil, fmt.Errorf("qjs new: %w", err)
	}
	defer rt.Close()
	ctx := rt.Context()

	injectFetch(ctx, q.transport)

	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	js := fmt.Sprintf(`await (%s)(%s)`, s.Body, string(argsJSON))

	result, err := ctx.Eval("script.js", qjs.Code(js), qjs.FlagAsync())
	if err != nil {
		return nil, fmt.Errorf("qjs eval: %w", err)
	}
	defer result.Free()

	jsonStr := stringify(ctx, result)
	if jsonStr == "undefined" || jsonStr == "" {
		return nil, fmt.Errorf("qjs returned empty result")
	}

	var out any
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		return nil, fmt.Errorf("unmarshal result: %w", err)
	}

	return out, nil
}

// stringify converts a QJS value to a JSON string via JSON.stringify().
func stringify(ctx *qjs.Context, val *qjs.Value) string {
	result, err := ctx.Global().GetPropertyStr("JSON").InvokeJS("stringify", val)
	if err != nil {
		return "{}"
	}
	defer result.Free()
	return result.String()
}

// injectFetch adds a fetch() function backed by the shared transport's HTTP client.
func injectFetch(ctx *qjs.Context, tp *transport.Transport) {
	ctx.SetAsyncFunc("fetch", func(this *qjs.This) {
		c := this.Context()

		url := ""
		method := "GET"
		var headers map[string]string
		var body string

		args := this.Args()
		if len(args) > 0 {
			url = args[0].String()
		}
		if len(args) > 1 && !args[1].IsUndefined() && !args[1].IsNull() {
			opts := args[1]
			m := opts.GetPropertyStr("method")
			if !m.IsUndefined() && !m.IsNull() {
				method = strings.ToUpper(m.String())
			}
			h := opts.GetPropertyStr("headers")
			if !h.IsUndefined() && !h.IsNull() {
				headersJSON := stringify(c, h)
				_ = json.Unmarshal([]byte(headersJSON), &headers)
			}
			b := opts.GetPropertyStr("body")
			if !b.IsUndefined() && !b.IsNull() {
				body = b.String()
			}
		}

		go func() {
			var bodyReader io.Reader
			if body != "" {
				bodyReader = strings.NewReader(body)
			}

			req, err := http.NewRequest(method, url, bodyReader)
			if err != nil {
				this.Promise().Reject(c.NewError(errors.New(err.Error())))
				return
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
			req.Header.Set("Accept", "application/json, text/plain, */*")
			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err := tp.Do(context.Background(), req)
			if err != nil {
				this.Promise().Reject(c.NewError(errors.New(err.Error())))
				return
			}
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				this.Promise().Reject(c.NewError(errors.New(err.Error())))
				return
			}

			respObj := c.NewObject()
			respObj.SetPropertyStr("ok", c.NewBool(resp.StatusCode >= 200 && resp.StatusCode < 300))
			respObj.SetPropertyStr("status", c.NewInt32(int32(resp.StatusCode)))
			respObj.SetPropertyStr("statusText", c.NewString(resp.Status))
			respObj.SetPropertyStr("_body", c.NewString(string(respBody)))

			this.Promise().Resolve(respObj)
		}()
	})

	ctx.Eval("fetch-polyfill.js", qjs.Code(`
		const _rawFetch = fetch;
		globalThis.fetch = async function(url, opts) {
			const resp = await _rawFetch(url, opts);
			return {
				ok: resp.ok,
				status: resp.status,
				statusText: resp.statusText,
				text: async () => resp._body,
				json: async () => JSON.parse(resp._body),
			};
		};
	`), qjs.FlagAsync())
}
