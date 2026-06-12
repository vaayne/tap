package browser

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

// TestEmulationSettingsIsEmpty verifies the IsEmpty sentinel covers all fields.
func TestEmulationSettingsIsEmpty(t *testing.T) {
	nilE := (*EmulationSettings)(nil)
	if !nilE.IsEmpty() {
		t.Fatal("nil EmulationSettings.IsEmpty() should be true")
	}
	zeroE := &EmulationSettings{}
	if !zeroE.IsEmpty() {
		t.Fatal("zero EmulationSettings.IsEmpty() should be true")
	}

	lat := 37.7749
	lng := -122.4194
	offline := true
	cases := []struct {
		name string
		e    *EmulationSettings
	}{
		{"viewport", &EmulationSettings{ViewportWidth: 1280, ViewportHeight: 720}},
		{"viewport scale", &EmulationSettings{ViewportWidth: 390, ViewportHeight: 844, ViewportScale: 3.0}},
		{"device", &EmulationSettings{DeviceName: "iPhone 14"}},
		{"geo", &EmulationSettings{GeoLat: &lat, GeoLng: &lng}},
		{"offline", &EmulationSettings{Offline: &offline}},
		{"headers", &EmulationSettings{Headers: map[string]string{"X-Foo": "bar"}}},
		{"media", &EmulationSettings{MediaScheme: "dark"}},
		{"useragent", &EmulationSettings{UserAgent: "TestBot/1.0"}},
	}
	for _, tc := range cases {
		if tc.e.IsEmpty() {
			t.Errorf("case %q: IsEmpty() should be false", tc.name)
		}
	}
}

// TestEmulationSettingsJSONRoundTrip verifies the struct serializes and
// deserializes correctly through the Store's JSON pipeline.
func TestEmulationSettingsJSONRoundTrip(t *testing.T) {
	lat := 37.7749
	lng := -122.4194
	offline := true

	orig := EmulationSettings{
		ViewportWidth:  1280,
		ViewportHeight: 720,
		ViewportScale:  2.0,
		DeviceName:     "Pixel 5",
		GeoLat:         &lat,
		GeoLng:         &lng,
		Offline:        &offline,
		Headers:        map[string]string{"Authorization": "Bearer tok", "X-Custom": "v"},
		MediaScheme:    "dark",
		UserAgent:      "TestBot/1.0",
	}

	data, err := json.Marshal(&orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got EmulationSettings
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ViewportWidth != orig.ViewportWidth {
		t.Errorf("ViewportWidth: got %d, want %d", got.ViewportWidth, orig.ViewportWidth)
	}
	if got.ViewportHeight != orig.ViewportHeight {
		t.Errorf("ViewportHeight: got %d, want %d", got.ViewportHeight, orig.ViewportHeight)
	}
	if got.ViewportScale != orig.ViewportScale {
		t.Errorf("ViewportScale: got %f, want %f", got.ViewportScale, orig.ViewportScale)
	}
	if got.DeviceName != orig.DeviceName {
		t.Errorf("DeviceName: got %q, want %q", got.DeviceName, orig.DeviceName)
	}
	if got.GeoLat == nil || *got.GeoLat != lat {
		t.Errorf("GeoLat: got %v, want %v", got.GeoLat, lat)
	}
	if got.GeoLng == nil || *got.GeoLng != lng {
		t.Errorf("GeoLng: got %v, want %v", got.GeoLng, lng)
	}
	if got.Offline == nil || *got.Offline != offline {
		t.Errorf("Offline: got %v, want %v", got.Offline, offline)
	}
	if len(got.Headers) != len(orig.Headers) {
		t.Errorf("Headers len: got %d, want %d", len(got.Headers), len(orig.Headers))
	}
	for k, v := range orig.Headers {
		if got.Headers[k] != v {
			t.Errorf("Headers[%q]: got %q, want %q", k, got.Headers[k], v)
		}
	}
	if got.MediaScheme != orig.MediaScheme {
		t.Errorf("MediaScheme: got %q, want %q", got.MediaScheme, orig.MediaScheme)
	}
	if got.UserAgent != orig.UserAgent {
		t.Errorf("UserAgent: got %q, want %q", got.UserAgent, orig.UserAgent)
	}
}

// TestEmulationSettingsStorePersistence verifies that EmulationSettings stored
// in a TabRecord survive a full Store save/load cycle.
func TestEmulationSettingsStorePersistence(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	session, err := NewLocalSession(DefaultSessionName, filepath.Join(t.TempDir(), "profiles", "default"), true, now)
	if err != nil {
		t.Fatalf("NewLocalSession: %v", err)
	}

	tab, err := NewTab("main", "target-123", "https://example.com", now)
	if err != nil {
		t.Fatalf("NewTab: %v", err)
	}

	lat := 48.8566
	lng := 2.3522
	offline := false
	tab.Emulation = &EmulationSettings{
		ViewportWidth:  390,
		ViewportHeight: 844,
		ViewportScale:  3.0,
		GeoLat:         &lat,
		GeoLng:         &lng,
		Offline:        &offline,
		Headers:        map[string]string{"X-Test": "1"},
		MediaScheme:    "light",
		UserAgent:      "TestAgent/2.0",
	}

	if err := store.Update(func(state *State) error {
		if err := state.CreateSession(session); err != nil {
			return err
		}
		return state.UpsertTab(DefaultSessionName, tab)
	}); err != nil {
		t.Fatalf("store.Update: %v", err)
	}

	// Reload from disk.
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("store.Load: %v", err)
	}

	s, err := loaded.ResolveSession(DefaultSessionName)
	if err != nil {
		t.Fatalf("ResolveSession: %v", err)
	}
	tt, ok := s.Tabs["main"]
	if !ok {
		t.Fatal("tab 'main' not found after reload")
	}
	if tt.Emulation == nil {
		t.Fatal("Emulation is nil after reload")
	}
	e := tt.Emulation
	if e.ViewportWidth != 390 || e.ViewportHeight != 844 {
		t.Errorf("viewport: got %dx%d, want 390x844", e.ViewportWidth, e.ViewportHeight)
	}
	if e.ViewportScale != 3.0 {
		t.Errorf("scale: got %f, want 3.0", e.ViewportScale)
	}
	if e.GeoLat == nil || *e.GeoLat != lat {
		t.Errorf("GeoLat: got %v, want %v", e.GeoLat, lat)
	}
	if e.Offline == nil || *e.Offline != false {
		t.Errorf("Offline: got %v, want false", e.Offline)
	}
	if e.Headers["X-Test"] != "1" {
		t.Errorf("Headers[X-Test]: got %q, want %q", e.Headers["X-Test"], "1")
	}
	if e.MediaScheme != "light" {
		t.Errorf("MediaScheme: got %q, want light", e.MediaScheme)
	}
	if e.UserAgent != "TestAgent/2.0" {
		t.Errorf("UserAgent: got %q, want TestAgent/2.0", e.UserAgent)
	}
}

// TestMergeEmulation verifies that mergeEmulation correctly overlays delta fields.
func TestMergeEmulation(t *testing.T) {
	dst := &EmulationSettings{
		ViewportWidth:  1920,
		ViewportHeight: 1080,
		MediaScheme:    "light",
	}
	offline := true
	src := &EmulationSettings{
		MediaScheme: "dark",
		Offline:     &offline,
		UserAgent:   "NewBot/1.0",
	}

	mergeEmulation(dst, src)

	// Untouched fields should persist.
	if dst.ViewportWidth != 1920 {
		t.Errorf("ViewportWidth changed: got %d", dst.ViewportWidth)
	}
	// Changed fields from src should be applied.
	if dst.MediaScheme != "dark" {
		t.Errorf("MediaScheme: got %q, want dark", dst.MediaScheme)
	}
	if dst.Offline == nil || !*dst.Offline {
		t.Errorf("Offline: expected true, got %v", dst.Offline)
	}
	if dst.UserAgent != "NewBot/1.0" {
		t.Errorf("UserAgent: got %q, want NewBot/1.0", dst.UserAgent)
	}
}

// TestMergeEmulationDevicePresetClearsViewport verifies that setting a device
// preset clears independent viewport/UA overrides to avoid stale re-application.
func TestMergeEmulationDevicePresetClearsViewport(t *testing.T) {
	dst := &EmulationSettings{
		ViewportWidth:  1280,
		ViewportHeight: 720,
		ViewportScale:  2.0,
		UserAgent:      "OldBot/1.0",
	}
	src := &EmulationSettings{DeviceName: "iPhone 14"}

	mergeEmulation(dst, src)

	if dst.DeviceName != "iPhone 14" {
		t.Errorf("DeviceName: got %q, want iPhone 14", dst.DeviceName)
	}
	// Device preset replaces independent viewport fields.
	if dst.ViewportWidth != 0 || dst.ViewportHeight != 0 || dst.ViewportScale != 0 {
		t.Errorf("viewport should be zeroed after device preset, got %dx%d scale=%f",
			dst.ViewportWidth, dst.ViewportHeight, dst.ViewportScale)
	}
	if dst.UserAgent != "" {
		t.Errorf("UserAgent should be cleared after device preset, got %q", dst.UserAgent)
	}
}

// TestLookupDevice verifies case-insensitive device lookup.
func TestLookupDevice(t *testing.T) {
	cases := []string{
		"iPhone 14",
		"iphone 14",
		"IPHONE 14",
		"Pixel 5",
		"pixel 5",
	}
	for _, name := range cases {
		d, err := lookupDevice(name)
		if err != nil {
			t.Errorf("lookupDevice(%q): unexpected error: %v", name, err)
			continue
		}
		if d.Name == "" {
			t.Errorf("lookupDevice(%q): returned empty name", name)
		}
	}

	// Unknown device should return error.
	_, err := lookupDevice("NonExistentDevice9999")
	if err == nil {
		t.Error("lookupDevice(unknown): expected error, got nil")
	}
}

// TestHeadersJSONParsing verifies that the JSON header parsing used by the CLI
// correctly handles standard and edge-case inputs.
func TestHeadersJSONParsing(t *testing.T) {
	cases := []struct {
		input   string
		wantLen int
		wantErr bool
	}{
		{`{"Authorization":"Bearer tok"}`, 1, false},
		{`{"X-A":"1","X-B":"2"}`, 2, false},
		{`{}`, 0, false},
		{`not-json`, 0, true},
		{`["array"]`, 0, true},
	}
	for _, tc := range cases {
		var headers map[string]string
		err := json.Unmarshal([]byte(tc.input), &headers)
		hasErr := err != nil
		if hasErr != tc.wantErr {
			t.Errorf("input %q: wantErr=%v gotErr=%v", tc.input, tc.wantErr, hasErr)
			continue
		}
		if !tc.wantErr && len(headers) != tc.wantLen {
			t.Errorf("input %q: len(headers)=%d, want %d", tc.input, len(headers), tc.wantLen)
		}
	}
}
