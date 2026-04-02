package browser

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestDebugURLToHTTP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "ws to http",
			input: "ws://127.0.0.1:9222/devtools/browser/abc",
			want:  "http://127.0.0.1:9222",
		},
		{
			name:  "wss to https",
			input: "wss://remote:9222/devtools/browser/abc",
			want:  "https://remote:9222",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid scheme http",
			input:   "http://127.0.0.1:9222/devtools/browser/abc",
			wantErr: true,
		},
		{
			name:    "invalid scheme https",
			input:   "https://127.0.0.1:9222/devtools/browser/abc",
			wantErr: true,
		},
		{
			name:  "strips path",
			input: "ws://localhost:9222/deep/path/here",
			want:  "http://localhost:9222",
		},
		{
			name:  "strips query",
			input: "ws://localhost:9222/devtools?token=abc",
			want:  "http://localhost:9222",
		},
		{
			name:  "strips fragment",
			input: "ws://localhost:9222/devtools#section",
			want:  "http://localhost:9222",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := debugURLToHTTP(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("debugURLToHTTP(%q) = %q, want error", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("debugURLToHTTP(%q) error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("debugURLToHTTP(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDebugURL(t *testing.T) {
	t.Run("devtools line present", func(t *testing.T) {
		input := "DevTools listening on ws://127.0.0.1:9222/devtools/browser/abc\n"
		got, err := parseDebugURL(strings.NewReader(input), 5*time.Second)
		if err != nil {
			t.Fatalf("parseDebugURL error: %v", err)
		}
		want := "ws://127.0.0.1:9222/devtools/browser/abc"
		if got != want {
			t.Fatalf("parseDebugURL = %q, want %q", got, want)
		}
	})

	t.Run("devtools line after noise", func(t *testing.T) {
		input := "Some startup noise\nAnother line\nDevTools listening on ws://127.0.0.1:5555/devtools/browser/xyz\n"
		got, err := parseDebugURL(strings.NewReader(input), 5*time.Second)
		if err != nil {
			t.Fatalf("parseDebugURL error: %v", err)
		}
		want := "ws://127.0.0.1:5555/devtools/browser/xyz"
		if got != want {
			t.Fatalf("parseDebugURL = %q, want %q", got, want)
		}
	})

	t.Run("pipe closed without line", func(t *testing.T) {
		input := "some output\nno devtools line here\n"
		_, err := parseDebugURL(strings.NewReader(input), 5*time.Second)
		if err == nil {
			t.Fatal("parseDebugURL should return error when pipe closes without DevTools line")
		}
	})

	t.Run("timeout with no output", func(t *testing.T) {
		// Use a pipe that blocks forever (never writes).
		r, w := io.Pipe()
		defer func() { _ = w.Close() }()

		_, err := parseDebugURL(r, 50*time.Millisecond)
		if err == nil {
			t.Fatal("parseDebugURL should return timeout error")
		}
		if !strings.Contains(err.Error(), "timed out") {
			t.Fatalf("parseDebugURL error = %q, want timeout error", err.Error())
		}
	})
}

func TestFindChrome(t *testing.T) {
	// Just verify findChrome doesn't panic. It may return a path or an error
	// depending on the environment.
	_, _ = findChrome()
}

func TestKillProcessNilRecord(t *testing.T) {
	if err := KillProcess(nil); err != nil {
		t.Fatalf("KillProcess(nil) = %v, want nil", err)
	}
}

func TestKillProcessZeroPID(t *testing.T) {
	if err := KillProcess(&ProcessRecord{PID: 0}); err != nil {
		t.Fatalf("KillProcess(PID=0) = %v, want nil", err)
	}
}

func TestCheckProcessNilRecord(t *testing.T) {
	err := CheckProcess(nil)
	if err == nil {
		t.Fatal("CheckProcess(nil) should return error")
	}
}
