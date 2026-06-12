package browser

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// EmulationSettings is the durable emulation configuration for a tab.
// All fields are optional: nil/zero means "leave untouched".
type EmulationSettings struct {
	// Viewport: width and height must both be non-zero to take effect; scale defaults to 1.
	ViewportWidth  int64   `json:"viewport_width,omitempty"`
	ViewportHeight int64   `json:"viewport_height,omitempty"`
	ViewportScale  float64 `json:"viewport_scale,omitempty"`
	// DeviceName is a chromedp/device preset name (e.g. "iPhone 14").
	DeviceName string `json:"device_name,omitempty"`
	// Geo override — both must be non-nil to take effect.
	GeoLat *float64 `json:"geo_lat,omitempty"`
	GeoLng *float64 `json:"geo_lng,omitempty"`
	// Offline sets navigator.onLine to false when true.
	Offline *bool `json:"offline,omitempty"`
	// Headers adds extra HTTP request headers.
	Headers map[string]string `json:"headers,omitempty"`
	// MediaScheme is "dark" or "light" for prefers-color-scheme.
	MediaScheme string `json:"media_scheme,omitempty"`
	// UserAgent overrides navigator.userAgent.
	UserAgent string `json:"user_agent,omitempty"`
}

// IsEmpty returns true when no field is set (nothing to apply).
func (e *EmulationSettings) IsEmpty() bool {
	if e == nil {
		return true
	}
	return e.ViewportWidth == 0 &&
		e.ViewportHeight == 0 &&
		e.ViewportScale == 0 &&
		e.DeviceName == "" &&
		e.GeoLat == nil &&
		e.GeoLng == nil &&
		e.Offline == nil &&
		len(e.Headers) == 0 &&
		e.MediaScheme == "" &&
		e.UserAgent == ""
}

// deviceCatalog maps lowercase device names to chromedp device.Info.
// Built from the full set of exported chromedp/device constants.
var deviceCatalog map[string]device.Info

func init() {
	entries := []device.Info{
		device.BlackberryPlayBook.Device(),
		device.BlackberryPlayBooklandscape.Device(),
		device.BlackBerryZ30.Device(),
		device.BlackBerryZ30landscape.Device(),
		device.GalaxyNote3.Device(),
		device.GalaxyNote3landscape.Device(),
		device.GalaxyNoteII.Device(),
		device.GalaxyNoteIIlandscape.Device(),
		device.GalaxySIII.Device(),
		device.GalaxySIIIlandscape.Device(),
		device.GalaxyS5.Device(),
		device.GalaxyS5landscape.Device(),
		device.GalaxyS8.Device(),
		device.GalaxyS8landscape.Device(),
		device.GalaxyS9.Device(),
		device.GalaxyS9landscape.Device(),
		device.GalaxyTabS4.Device(),
		device.GalaxyTabS4landscape.Device(),
		device.IPad.Device(),
		device.IPadlandscape.Device(),
		device.IPadgen7.Device(),
		device.IPadgen7landscape.Device(),
		device.IPadMini.Device(),
		device.IPadMinilandscape.Device(),
		device.IPadPro.Device(),
		device.IPadProlandscape.Device(),
		device.IPadPro11.Device(),
		device.IPadPro11landscape.Device(),
		device.IPhone4.Device(),
		device.IPhone4landscape.Device(),
		device.IPhone5.Device(),
		device.IPhone5landscape.Device(),
		device.IPhone6.Device(),
		device.IPhone6landscape.Device(),
		device.IPhone6Plus.Device(),
		device.IPhone6Pluslandscape.Device(),
		device.IPhone7.Device(),
		device.IPhone7landscape.Device(),
		device.IPhone7Plus.Device(),
		device.IPhone7Pluslandscape.Device(),
		device.IPhone8.Device(),
		device.IPhone8landscape.Device(),
		device.IPhone8Plus.Device(),
		device.IPhone8Pluslandscape.Device(),
		device.IPhoneX.Device(),
		device.IPhoneXlandscape.Device(),
		device.IPhoneXR.Device(),
		device.IPhoneXRlandscape.Device(),
		device.IPhoneSE.Device(),
		device.IPhoneSElandscape.Device(),
		device.IPhone12Pro.Device(),
		device.IPhone12Prolandscape.Device(),
		device.IPhone14.Device(),
		device.IPhone14landscape.Device(),
		device.IPhone14Plus.Device(),
		device.IPhone14Pluslandscape.Device(),
		device.IPhone14Pro.Device(),
		device.IPhone14Prolandscape.Device(),
		device.IPhone14ProMax.Device(),
		device.IPhone14ProMaxlandscape.Device(),
		device.JioPhone2.Device(),
		device.JioPhone2landscape.Device(),
		device.KindleFireHDX.Device(),
		device.KindleFireHDXlandscape.Device(),
		device.LGOptimusL70.Device(),
		device.LGOptimusL70landscape.Device(),
		device.MicrosoftLumia550.Device(),
		device.MicrosoftLumia950.Device(),
		device.MicrosoftLumia950landscape.Device(),
		device.Nexus10.Device(),
		device.Nexus10landscape.Device(),
		device.Nexus4.Device(),
		device.Nexus4landscape.Device(),
		device.Nexus5.Device(),
		device.Nexus5landscape.Device(),
		device.Nexus5X.Device(),
		device.Nexus5Xlandscape.Device(),
		device.Nexus6.Device(),
		device.Nexus6landscape.Device(),
		device.Nexus6P.Device(),
		device.Nexus6Plandscape.Device(),
		device.Nexus7.Device(),
		device.Nexus7landscape.Device(),
		device.NokiaLumia520.Device(),
		device.NokiaLumia520landscape.Device(),
		device.NokiaN9.Device(),
		device.NokiaN9landscape.Device(),
		device.Pixel2.Device(),
		device.Pixel2landscape.Device(),
		device.Pixel2XL.Device(),
		device.Pixel2XLlandscape.Device(),
		device.Pixel3.Device(),
		device.Pixel3landscape.Device(),
		device.Pixel4.Device(),
		device.Pixel4landscape.Device(),
		device.Pixel4a5G.Device(),
		device.Pixel4a5Glandscape.Device(),
		device.Pixel5.Device(),
		device.Pixel5landscape.Device(),
	}

	deviceCatalog = make(map[string]device.Info, len(entries))
	for _, d := range entries {
		if d.Name != "" {
			deviceCatalog[strings.ToLower(d.Name)] = d
		}
	}
}

// lookupDevice resolves a human-readable device name to a chromedp device.Info.
// Matching is case-insensitive.
func lookupDevice(name string) (device.Info, error) {
	d, ok := deviceCatalog[strings.ToLower(strings.TrimSpace(name))]
	if !ok {
		return device.Info{}, fmt.Errorf("unknown device %q; use the exact name from chromedp/device (e.g. \"iPhone 14\")", name)
	}
	return d, nil
}

// ApplyEmulation returns a chromedp.Action that applies all non-zero fields of e
// to the currently attached CDP target. Safe to call with a nil or empty e (no-op).
func ApplyEmulation(e *EmulationSettings) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if e.IsEmpty() {
			return nil
		}

		// Device preset sets viewport + UA + touch in one shot.
		if e.DeviceName != "" {
			d, err := lookupDevice(e.DeviceName)
			if err != nil {
				return err
			}
			if err := chromedp.Emulate(d).Do(ctx); err != nil {
				return fmt.Errorf("apply device %q: %w", e.DeviceName, err)
			}
		}

		// Explicit viewport overrides device preset when both are set.
		if e.ViewportWidth > 0 && e.ViewportHeight > 0 {
			scale := e.ViewportScale
			if scale == 0 {
				scale = 1.0
			}
			if err := emulation.SetDeviceMetricsOverride(e.ViewportWidth, e.ViewportHeight, scale, false).Do(ctx); err != nil {
				return fmt.Errorf("set viewport: %w", err)
			}
		}

		// Explicit UA overrides device preset when both are set.
		if e.UserAgent != "" {
			if err := emulation.SetUserAgentOverride(e.UserAgent).Do(ctx); err != nil {
				return fmt.Errorf("set user agent: %w", err)
			}
		}

		// Geolocation override — both lat and lng must be present.
		if e.GeoLat != nil && e.GeoLng != nil {
			p := emulation.SetGeolocationOverride().
				WithLatitude(*e.GeoLat).
				WithLongitude(*e.GeoLng)
			if err := p.Do(ctx); err != nil {
				return fmt.Errorf("set geo: %w", err)
			}
		}

		// Offline / online mode.
		if e.Offline != nil {
			// OverrideNetworkState requires the Network domain to be enabled.
			if err := network.Enable().Do(ctx); err != nil {
				return fmt.Errorf("enable network domain: %w", err)
			}
			// latency=0, down/up=-1 means "no constraint" (online); down/up=0 means offline.
			latency := 0.0
			down := -1.0
			up := -1.0
			if *e.Offline {
				down = 0
				up = 0
			}
			if err := network.OverrideNetworkState(*e.Offline, latency, down, up).Do(ctx); err != nil {
				return fmt.Errorf("set offline mode: %w", err)
			}
		}

		// Extra HTTP headers.
		if len(e.Headers) > 0 {
			h := make(network.Headers, len(e.Headers))
			for k, v := range e.Headers {
				h[k] = v
			}
			if err := network.Enable().Do(ctx); err != nil {
				return fmt.Errorf("enable network domain: %w", err)
			}
			if err := network.SetExtraHTTPHeaders(h).Do(ctx); err != nil {
				return fmt.Errorf("set extra http headers: %w", err)
			}
		}

		// Prefers-color-scheme media emulation.
		if e.MediaScheme != "" {
			feat := &emulation.MediaFeature{
				Name:  "prefers-color-scheme",
				Value: e.MediaScheme,
			}
			p := emulation.SetEmulatedMedia().WithFeatures([]*emulation.MediaFeature{feat})
			if err := p.Do(ctx); err != nil {
				return fmt.Errorf("set media scheme: %w", err)
			}
		}

		return nil
	})
}

// ApplyEmulationTarget applies persisted EmulationSettings to a live CDP target.
// It is a no-op when e is nil or empty.
func ApplyEmulationTarget(ctx context.Context, debugURL, targetID string, e *EmulationSettings) error {
	if e.IsEmpty() {
		return nil
	}
	return withTarget(ctx, debugURL, targetID, ApplyEmulation(e))
}
