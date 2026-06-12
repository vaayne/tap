package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// ---------------------------------------------------------------------------
// Web storage (localStorage / sessionStorage) via Runtime.evaluate
// ---------------------------------------------------------------------------

// storageGetAllJS returns JS that reads all entries from a named storage object.
func storageGetAllJS(storeName string) string {
	return fmt.Sprintf(`(() => {
  const s = window[%q];
  const out = {};
  for (let i = 0; i < s.length; i++) {
    const k = s.key(i);
    out[k] = s.getItem(k);
  }
  return JSON.stringify(out);
})()`, storeName)
}

func storageGetKeyJS(storeName, key string) string {
	return fmt.Sprintf(`(() => {
  return window[%q].getItem(%q);
})()`, storeName, key)
}

func storageSetJS(storeName, key, value string) string {
	return fmt.Sprintf(`(() => {
  window[%q].setItem(%q, %q);
  return true;
})()`, storeName, key, value)
}

func storageClearJS(storeName string) string {
	return fmt.Sprintf(`(() => {
  window[%q].clear();
  return true;
})()`, storeName)
}

// evalOpts returns a standard set of Evaluate options for storage JS (return by value, no await).
func evalOpts(p *runtime.EvaluateParams) *runtime.EvaluateParams {
	return p.WithReturnByValue(true)
}

// GetStorageAll returns all key/value entries from the named storage ("localStorage" or "sessionStorage").
func GetStorageAll(ctx context.Context, debugURL, targetID, storeName string) (map[string]string, error) {
	var raw string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(storageGetAllJS(storeName), &raw, evalOpts),
	)
	if err != nil {
		return nil, fmt.Errorf("storage get-all %s: %w", storeName, err)
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, fmt.Errorf("storage get-all %s: parse result: %w", storeName, err)
	}
	return out, nil
}

// GetStorageKey returns the value for a single key from the named storage.
// Returns ("", nil) when the key does not exist (localStorage.getItem returns null).
func GetStorageKey(ctx context.Context, debugURL, targetID, storeName, key string) (string, error) {
	var val any
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(storageGetKeyJS(storeName, key), &val, evalOpts),
	)
	if err != nil {
		return "", fmt.Errorf("storage get %s[%q]: %w", storeName, key, err)
	}
	if val == nil {
		return "", nil
	}
	s, ok := val.(string)
	if !ok {
		return fmt.Sprintf("%v", val), nil
	}
	return s, nil
}

// SetStorageKey sets a key in the named storage.
func SetStorageKey(ctx context.Context, debugURL, targetID, storeName, key, value string) error {
	var ok bool
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(storageSetJS(storeName, key, value), &ok, evalOpts),
	)
	if err != nil {
		return fmt.Errorf("storage set %s[%q]: %w", storeName, key, err)
	}
	return nil
}

// ClearStorage clears all entries from the named storage.
func ClearStorage(ctx context.Context, debugURL, targetID, storeName string) error {
	var ok bool
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(storageClearJS(storeName), &ok, evalOpts),
	)
	if err != nil {
		return fmt.Errorf("storage clear %s: %w", storeName, err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// State export/import (cookies + localStorage)
// ---------------------------------------------------------------------------

// currentOriginJS returns the origin of the current page.
const currentOriginJS = `location.origin`

// GetCurrentOrigin returns the window.location.origin of the current page.
func GetCurrentOrigin(ctx context.Context, debugURL, targetID string) (string, error) {
	var origin string
	err := withTarget(ctx, debugURL, targetID,
		chromedp.Evaluate(currentOriginJS, &origin, evalOpts),
	)
	if err != nil {
		return "", fmt.Errorf("get current origin: %w", err)
	}
	return origin, nil
}

// GetAllCookiesTarget returns all cookies visible in the current page context.
// CDP's network.GetCookies() returns cookies for the page's URL/domain context.
// For a comprehensive export we also fall back to all cookies the browser holds;
// the cdproto version used here does not expose GetAllCookies so we use GetCookies.
func GetAllCookiesTarget(ctx context.Context, debugURL, targetID string) ([]StateCookie, error) {
	var cookies []StateCookie
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			raw, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, c := range raw {
				cookies = append(cookies, StateCookie{
					Name:     c.Name,
					Value:    c.Value,
					Domain:   c.Domain,
					Path:     c.Path,
					Expires:  c.Expires,
					HTTPOnly: c.HTTPOnly,
					Secure:   c.Secure,
					SameSite: string(c.SameSite),
				})
			}
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("get all cookies: %w", err)
	}
	return cookies, nil
}

// SetAllCookiesTarget sets a batch of cookies with full attribute support.
// This is separate from SetCookieTarget (which only supports name/value/domain/path)
// to preserve expires/httpOnly/secure/sameSite when restoring saved state.
func SetAllCookiesTarget(ctx context.Context, debugURL, targetID string, cookies []StateCookie) error {
	err := withTarget(ctx, debugURL, targetID,
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, c := range cookies {
				p := network.SetCookie(c.Name, c.Value)
				if c.Domain != "" {
					p = p.WithDomain(c.Domain)
				}
				if c.Path != "" {
					p = p.WithPath(c.Path)
				}
				if c.Expires != 0 {
					ts := cdp.TimeSinceEpoch(time.Unix(int64(c.Expires), 0))
					p = p.WithExpires(&ts)
				}
				if c.HTTPOnly {
					p = p.WithHTTPOnly(c.HTTPOnly)
				}
				if c.Secure {
					p = p.WithSecure(c.Secure)
				}
				if c.SameSite != "" {
					p = p.WithSameSite(network.CookieSameSite(c.SameSite))
				}
				if err := p.Do(ctx); err != nil {
					return fmt.Errorf("set cookie %q: %w", c.Name, err)
				}
			}
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("set all cookies: %w", err)
	}
	return nil
}
