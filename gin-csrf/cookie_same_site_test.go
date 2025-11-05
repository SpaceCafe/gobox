package csrf_test

import (
	"net/http"
	"testing"

	csrf "github.com/spacecafe/gobox/gin-csrf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sameSite csrf.CookieSameSite
		want     string
	}{
		{"LaxMode", csrf.CookieSameSite{http.SameSiteLaxMode}, "lax"},
		{"StrictMode", csrf.CookieSameSite{http.SameSiteStrictMode}, "strict"},
		{"NoneMode", csrf.CookieSameSite{http.SameSiteNoneMode}, "none"},
		{"Default", csrf.CookieSameSite{http.SameSiteDefaultMode}, ""},
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
		sameSite csrf.CookieSameSite
		want     string
	}{
		{"LaxMode", csrf.CookieSameSite{http.SameSiteLaxMode}, "lax"},
		{"StrictMode", csrf.CookieSameSite{http.SameSiteStrictMode}, "strict"},
		{"NoneMode", csrf.CookieSameSite{http.SameSiteNoneMode}, "none"},
		{"Default", csrf.CookieSameSite{http.SameSiteDefaultMode}, ""},
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
		want    csrf.CookieSameSite
		wantErr bool
	}{
		{"LaxMode", "lax", csrf.CookieSameSite{http.SameSiteLaxMode}, false},
		{"StrictMode", "strict", csrf.CookieSameSite{http.SameSiteStrictMode}, false},
		{"NoneMode", "none", csrf.CookieSameSite{http.SameSiteNoneMode}, false},
		{"Invalid", "invalid", csrf.CookieSameSite{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var css csrf.CookieSameSite

			err := css.UnmarshalText([]byte(tt.input))

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, css)
		})
	}
}
