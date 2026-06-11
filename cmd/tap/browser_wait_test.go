package main

import (
	"testing"
	"time"

	"github.com/vaayne/tap/browser"
)

func TestParseDurationArg(t *testing.T) {
	tests := []struct {
		input  string
		want   time.Duration
		wantOK bool
	}{
		{"1000", 1000 * time.Millisecond, true},
		{"500", 500 * time.Millisecond, true},
		{"0", 0, true},
		{"2s", 2 * time.Second, true},
		{"1.5s", 1500 * time.Millisecond, true},
		{"300ms", 300 * time.Millisecond, true},
		{"#selector", 0, false},
		{".class", 0, false},
		{"div > span", 0, false},
		{"body", 0, false},
	}

	for _, tt := range tests {
		got, ok := parseDurationArg(tt.input)
		if ok != tt.wantOK {
			t.Errorf("parseDurationArg(%q): ok=%v want %v", tt.input, ok, tt.wantOK)
			continue
		}
		if ok && got != tt.want {
			t.Errorf("parseDurationArg(%q): duration=%v want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseElementState(t *testing.T) {
	tests := []struct {
		input   string
		want    browser.ElementState
		wantErr bool
	}{
		{"visible", browser.ElementVisible, false},
		{"", browser.ElementVisible, false},
		{"hidden", browser.ElementHidden, false},
		{"attached", browser.ElementAttached, false},
		{"detached", browser.ElementDetached, false},
		{"bogus", browser.ElementVisible, true},
	}

	for _, tt := range tests {
		got, err := parseElementState(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseElementState(%q): err=%v wantErr=%v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("parseElementState(%q): got %v want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseLoadState(t *testing.T) {
	tests := []struct {
		input   string
		want    browser.LoadState
		wantErr bool
	}{
		{"load", browser.LoadStateLoad, false},
		{"domcontentloaded", browser.LoadStateDOMContentLoaded, false},
		{"networkidle", browser.LoadStateNetworkIdle, false},
		{"ready", browser.LoadStateLoad, true},
		{"", browser.LoadStateLoad, true},
	}

	for _, tt := range tests {
		got, err := parseLoadState(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseLoadState(%q): err=%v wantErr=%v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("parseLoadState(%q): got %v want %v", tt.input, got, tt.want)
		}
	}
}
