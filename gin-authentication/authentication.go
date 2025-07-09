package authentication

import (
	"crypto/subtle"
	"errors"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("password is not valid")

	//nolint:gochecknoglobals // Maintain a set of predefined compareFuncs that are used throughout the application.
	compareFuncs = []struct {
		prefix  string
		compare func([]byte, []byte) error
	}{
		// Blank passwords.
		{"", compareBlankPasswords},

		// Bcrypt hashed passwords.
		// There are different bcrypt prefixes:
		// - "$2a$" is used by versions up to 1.0.4, has a known bug with handling of 8-bit characters.
		// - "$2x$" was added as a migration path for systems with "$2a$" prefix and still has a bug.
		// - "$2y$" is used in all later versions and should be used by modern systems.
		// - "$2b$" was introduced by OpenBSD 5.5., which behaves exactly like "$2y$".
		{"$2a$", bcrypt.CompareHashAndPassword},
		{"$2b$", bcrypt.CompareHashAndPassword},
		{"$2x$", bcrypt.CompareHashAndPassword},
		{"$2y$", bcrypt.CompareHashAndPassword},
	}
)

// New is a middleware function that handles HTTP Basic Auth, Bearer tokens and API Key header.
// If Basic Auth failed or not provided, it tries to get API key from the headers "API-Key" and "X-API-Key".
// The username in Basic Auth is ignored, but must contain at least one character.
func New(config *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenticated := false

		// Get API Key from the header if present.
		password := ctx.GetHeader(config.HeaderName)

		// Compare the provided password with server's configured API key using constant time comparison to prevent timing attacks.
		// If password matches, continue handling original request as user is authenticated.
		if password, ok := strings.CutPrefix(password, config.HeaderValuePrefix); ok {
			password = strings.TrimSpace(password)
			for i := range config.APIKeys {
				if ComparePasswords([]byte(config.APIKeys[i]), []byte(password)) {
					authenticated = true
				}
			}
		}

		// Retrieve the password from HTTP Basic Auth header if they exist.
		if username, password, ok := ctx.Request.BasicAuth(); ok {
			if _, ok = config.Users[username]; ok && ComparePasswords([]byte(config.Users[username]), []byte(password)) {
				authenticated = true
			}
		}

		// If the request is authenticated, continue processing.
		if authenticated {
			ctx.Next()
			return
		}

		// If password doesn't match or not provided at all, send a "401 Unauthorized" response with
		// "WWW-Authenticate" header to request client for valid credentials.
		ctx.Header("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		_ = ctx.Error(problems.ProblemUnauthorized)
		ctx.Abort()
	}
}

// ComparePasswords compares an inputted password with a hashed password using one of the functions in compareFuncs.
func ComparePasswords(hashedPassword, password []byte) bool {
	// Deny empty passwords.
	if len(hashedPassword) == 0 || len(password) == 0 {
		return false
	}

	// Match used hashing algorithm.
	compare := compareFuncs[0].compare
	for _, compareFunc := range compareFuncs[1:] {
		if strings.HasPrefix(string(hashedPassword), compareFunc.prefix) {
			compare = compareFunc.compare
			break
		}
	}
	return compare(hashedPassword, password) == nil
}

// compareBlankPasswords checks the provided passwords using constant time comparison to prevent timing attacks.
func compareBlankPasswords(hashedPassword, password []byte) error {
	lengthStored := len(hashedPassword)
	lengthUser := len(password)

	// Check for potential overflow and ensure the lengths are equal using constant time comparison.
	if lengthStored > math.MaxInt32 || lengthUser > math.MaxInt32 || subtle.ConstantTimeEq(int32(lengthStored), int32(lengthUser)) == 0 {

		// Perform a dummy comparison to mitigate timing attacks if the lengths do not match.
		subtle.ConstantTimeCompare(password, password)
		return ErrInvalidPassword
	}

	// Compare the hashed password with the user-provided password using constant time comparison.
	if subtle.ConstantTimeCompare(hashedPassword, password) == 1 {
		return nil
	}

	return ErrInvalidPassword
}
