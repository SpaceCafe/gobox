package authentication

import (
	"crypto/subtle"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	//nolint:gochecknoglobals // Maintain a set of predefined bcrypt prefixes that are used throughout the application.
	BcryptHashPrefixes = []string{"$2a$", "$2b$", "$2x$", "$2y$"}
)

// CompareSecrets compares two secrets (passwords or hashes) for equality.
// It handles both plain text and bcrypt hashed secrets.
func CompareSecrets(expected, actual string) error {
	compareFn := comparePasswords

	expectedBytes := []byte(expected)
	actualBytes := []byte(actual)

	for _, prefix := range BcryptHashPrefixes {
		if strings.HasPrefix(expected, prefix) {
			compareFn = bcrypt.CompareHashAndPassword
		}
	}

	return compareFn(expectedBytes, actualBytes)
}

// comparePasswords compares two passwords for equality.
// Its behavior is undefined if the password length is > 2**31-1.
func comparePasswords(expected, actual []byte) error {
	if subtle.ConstantTimeCompare(expected, actual) == 1 {
		return nil
	}
	return ErrSecretsNotEqual
}
