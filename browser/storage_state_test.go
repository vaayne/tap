package browser

import (
	"encoding/json"
	"testing"
)

// TestStorageStateRoundTrip validates that StorageState can be marshalled and
// unmarshalled without data loss — a basic sanity check for the state file format.
func TestStorageStateRoundTrip(t *testing.T) {
	original := StorageState{
		Cookies: []StateCookie{
			{
				Name:     "session",
				Value:    "abc123",
				Domain:   "example.com",
				Path:     "/",
				Expires:  1700000000,
				HTTPOnly: true,
				Secure:   true,
				SameSite: "Lax",
			},
			{
				Name:   "pref",
				Value:  "dark",
				Domain: ".example.com",
				Path:   "/",
			},
		},
		Origins: []OriginStorage{
			{
				Origin: "https://example.com",
				LocalStorage: []StorageEntry{
					{Name: "token", Value: "tok_xyz"},
					{Name: "user_id", Value: "42"},
				},
			},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded StorageState
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(decoded.Cookies) != len(original.Cookies) {
		t.Errorf("cookie count: got %d, want %d", len(decoded.Cookies), len(original.Cookies))
	}
	if len(decoded.Origins) != len(original.Origins) {
		t.Fatalf("origin count: got %d, want %d", len(decoded.Origins), len(original.Origins))
	}

	got := decoded.Cookies[0]
	want := original.Cookies[0]
	if got.Name != want.Name || got.Value != want.Value {
		t.Errorf("cookie[0]: got {%q,%q}, want {%q,%q}", got.Name, got.Value, want.Name, want.Value)
	}
	if got.HTTPOnly != want.HTTPOnly || got.Secure != want.Secure {
		t.Errorf("cookie[0] flags: got httpOnly=%v secure=%v, want httpOnly=%v secure=%v",
			got.HTTPOnly, got.Secure, want.HTTPOnly, want.Secure)
	}
	if got.SameSite != want.SameSite {
		t.Errorf("cookie[0] sameSite: got %q, want %q", got.SameSite, want.SameSite)
	}
	if got.Expires != want.Expires {
		t.Errorf("cookie[0] expires: got %v, want %v", got.Expires, want.Expires)
	}

	o := decoded.Origins[0]
	if o.Origin != original.Origins[0].Origin {
		t.Errorf("origin: got %q, want %q", o.Origin, original.Origins[0].Origin)
	}
	if len(o.LocalStorage) != len(original.Origins[0].LocalStorage) {
		t.Errorf("localStorage count: got %d, want %d", len(o.LocalStorage), len(original.Origins[0].LocalStorage))
	}
}

// TestStorageStatePlaywrightShape verifies the JSON keys match the Playwright
// storageState format that external tools (e.g. playwright CLI) expect.
func TestStorageStatePlaywrightShape(t *testing.T) {
	state := StorageState{
		Cookies: []StateCookie{
			{Name: "sid", Value: "v", Domain: "x.com", Path: "/", Expires: 0, HTTPOnly: false, Secure: false, SameSite: ""},
		},
		Origins: []OriginStorage{
			{Origin: "https://x.com", LocalStorage: []StorageEntry{{Name: "k", Value: "v"}}},
		},
	}
	data, _ := json.Marshal(state)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("outer parse: %v", err)
	}
	if _, ok := raw["cookies"]; !ok {
		t.Error("missing top-level \"cookies\" key")
	}
	if _, ok := raw["origins"]; !ok {
		t.Error("missing top-level \"origins\" key")
	}

	var cookies []map[string]json.RawMessage
	if err := json.Unmarshal(raw["cookies"], &cookies); err != nil {
		t.Fatalf("cookies parse: %v", err)
	}
	if len(cookies) == 0 {
		t.Fatal("no cookies in output")
	}
	for _, key := range []string{"name", "value", "domain", "path", "expires", "httpOnly", "secure"} {
		if _, ok := cookies[0][key]; !ok {
			t.Errorf("cookie missing key %q", key)
		}
	}

	var origins []map[string]json.RawMessage
	if err := json.Unmarshal(raw["origins"], &origins); err != nil {
		t.Fatalf("origins parse: %v", err)
	}
	if len(origins) == 0 {
		t.Fatal("no origins in output")
	}
	if _, ok := origins[0]["origin"]; !ok {
		t.Error("origin entry missing \"origin\" key")
	}
	if _, ok := origins[0]["localStorage"]; !ok {
		t.Error("origin entry missing \"localStorage\" key")
	}
}

// TestStorageStateEmptyState verifies that an empty state serialises cleanly.
func TestStorageStateEmptyState(t *testing.T) {
	state := StorageState{}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("marshal empty state: %v", err)
	}
	var decoded StorageState
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal empty state: %v", err)
	}
	if len(decoded.Cookies) != 0 {
		t.Errorf("expected 0 cookies, got %d", len(decoded.Cookies))
	}
	if len(decoded.Origins) != 0 {
		t.Errorf("expected 0 origins, got %d", len(decoded.Origins))
	}
}
