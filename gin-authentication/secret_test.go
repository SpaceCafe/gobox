package authentication_test

import (
	"testing"

	authentication "github.com/spacecafe/gobox/gin-authentication"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func hashBcryptSecret(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

func TestCompareSecrets(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected string
		actual   string
		wantErr  bool
	}{
		{"empty hashedPassword and password", "", "", false},
		{"empty hashedPassword", "", "secret", true},
		{"empty password", "secret", "", true},
		{"different blank passwords", "another secret", string("secret"), true},
		{"same blank passwords", "secret", "secret", false},
		{"same blank unicode password", "🔐secret", "🔐secret", false},
		{"different bcrypt hashed password", hashBcryptSecret("another secret"), "secret", true},
		{"same bcrypt hashed password", hashBcryptSecret("secret"), "secret", false},
		{"same bcrypt hashed unicode password", hashBcryptSecret("🔐secret"), "🔐secret", false},
		{
			"bcrypt '$2a$' hashed password",
			"$2a$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC",
			"secret",
			false,
		},
		{
			"bcrypt '$2b$' hashed password",
			"$2b$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC",
			"secret",
			false,
		},
		{
			"bcrypt '$2x$' hashed password",
			"$2x$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC",
			"secret",
			false,
		},
		{
			"bcrypt '$2y$' hashed password",
			"$2y$10$bbtYMfvxxXpRDTQEUv4EneR5figrz88R/j14RCbyxiNJweR4vBzkC",
			"secret",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := authentication.CompareSecrets(tt.expected, tt.actual)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
