package csrf

import (
	"bytes"
	"net/http"
)

// CookieSameSite represents the CookieSameSite attribute of a cookie, encapsulating http.SameSite.
type CookieSameSite struct {
	http.SameSite
}

// MarshalText converts the CookieSameSite value to a byte slice.
func (r *CookieSameSite) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// String returns the string representation of the CookieSameSite value.
func (r *CookieSameSite) String() string {
	switch r.SameSite {
	case http.SameSiteLaxMode:
		return "lax"
	case http.SameSiteStrictMode:
		return "strict"
	case http.SameSiteNoneMode:
		return "none"
	case http.SameSiteDefaultMode:
		return ""
	default:
		return ""
	}
}

// UnmarshalText converts the byte slice to a CookieSameSite value.
func (r *CookieSameSite) UnmarshalText(text []byte) error {
	switch string(bytes.ToLower(text)) {
	case "lax":
		r.SameSite = http.SameSiteLaxMode
	case "strict":
		r.SameSite = http.SameSiteStrictMode
	case "none":
		r.SameSite = http.SameSiteNoneMode
	default:
		return ErrInvalidCookieSameSite
	}

	return nil
}
