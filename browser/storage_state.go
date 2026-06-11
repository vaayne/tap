package browser

// StorageState is a Playwright storageState-compatible JSON structure for
// exporting and importing auth state (cookies + web storage).
type StorageState struct {
	Cookies []StateCookie   `json:"cookies"`
	Origins []OriginStorage `json:"origins"`
}

// StateCookie mirrors the Playwright cookie shape with all attributes needed
// for full-fidelity round-trips through CDP network.SetCookie.
type StateCookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	SameSite string  `json:"sameSite,omitempty"`
}

// OriginStorage holds the localStorage entries for a single origin.
type OriginStorage struct {
	Origin       string         `json:"origin"`
	LocalStorage []StorageEntry `json:"localStorage"`
}

// StorageEntry is a key/value pair from localStorage or sessionStorage.
type StorageEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
