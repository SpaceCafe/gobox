package jwt_test

import (
	"net/http"
	"testing"

	"github.com/spacecafe/gobox/gin-authentication/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sameSite jwt.CookieSameSite
		want     string
	}{
		{"LaxMode", jwt.CookieSameSite{http.SameSiteLaxMode}, "lax"},
		{"StrictMode", jwt.CookieSameSite{http.SameSiteStrictMode}, "strict"},
		{"NoneMode", jwt.CookieSameSite{http.SameSiteNoneMode}, "none"},
		{"Default", jwt.CookieSameSite{http.SameSiteDefaultMode}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.sameSite.String())
		})
	}
}

func TestMarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sameSite jwt.CookieSameSite
		want     string
	}{
		{"LaxMode", jwt.CookieSameSite{http.SameSiteLaxMode}, "lax"},
		{"StrictMode", jwt.CookieSameSite{http.SameSiteStrictMode}, "strict"},
		{"NoneMode", jwt.CookieSameSite{http.SameSiteNoneMode}, "none"},
		{"Default", jwt.CookieSameSite{http.SameSiteDefaultMode}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.sameSite.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestUnmarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    jwt.CookieSameSite
		wantErr bool
	}{
		{"LaxMode", "lax", jwt.CookieSameSite{http.SameSiteLaxMode}, false},
		{"StrictMode", "strict", jwt.CookieSameSite{http.SameSiteStrictMode}, false},
		{"NoneMode", "none", jwt.CookieSameSite{http.SameSiteNoneMode}, false},
		{"Invalid", "invalid", jwt.CookieSameSite{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var css jwt.CookieSameSite

			err := css.UnmarshalText([]byte(tt.input))

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, css)
		})
	}
}
