package jwt

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		sameSite CookieSameSite
		want     string
	}{
		{"LaxMode", CookieSameSite{http.SameSiteLaxMode}, "lax"},
		{"StrictMode", CookieSameSite{http.SameSiteStrictMode}, "strict"},
		{"NoneMode", CookieSameSite{http.SameSiteNoneMode}, "none"},
		{"Default", CookieSameSite{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sameSite.String())
		})
	}
}

func TestMarshalText(t *testing.T) {
	tests := []struct {
		name     string
		sameSite CookieSameSite
		want     string
	}{
		{"LaxMode", CookieSameSite{http.SameSiteLaxMode}, "lax"},
		{"StrictMode", CookieSameSite{http.SameSiteStrictMode}, "strict"},
		{"NoneMode", CookieSameSite{http.SameSiteNoneMode}, "none"},
		{"Default", CookieSameSite{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sameSite.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    CookieSameSite
		wantErr bool
	}{
		{"LaxMode", "lax", CookieSameSite{http.SameSiteLaxMode}, false},
		{"StrictMode", "strict", CookieSameSite{http.SameSiteStrictMode}, false},
		{"NoneMode", "none", CookieSameSite{http.SameSiteNoneMode}, false},
		{"Invalid", "invalid", CookieSameSite{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var css CookieSameSite
			err := css.UnmarshalText([]byte(tt.input))
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, css)
		})
	}
}
