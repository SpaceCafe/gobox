package csrf

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	authentication "github.com/spacecafe/gobox/gin-authentication"
)

// Token represents a CSRF token with associated configuration, context, session ID, and random bytes.
type Token struct {
	cfg       *Config
	ctx       *gin.Context
	sessionID string
	random    []byte
}

// NewToken generates a new CSRF token using the provided configuration.
func NewToken(cfg *Config, ctx *gin.Context) (*Token, error) {
	token, err := newToken(cfg, ctx)
	if err != nil {
		return nil, err
	}

	token.random = make([]byte, 64) //nolint:mnd // 64 bytes of entropy is sufficient.

	_, err = rand.Read(token.random)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return token, nil
}

// ValidateToken checks if the provided tokenString is valid for the given Config and gin.Context.
// It extracts the HMAC from the tokenString and compares it with the expected HMAC using subtle.ConstantTimeCompare.
// Returns an error if the token is invalid or any step fails.
func ValidateToken(cfg *Config, ctx *gin.Context, tokenString string) error {
	var hmacFromRequest []byte

	token, err := newToken(cfg, ctx)
	if err != nil {
		return err
	}

	hmacFromRequest, token.random, err = parseTokenString(tokenString)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(hmacFromRequest, token.HMAC()) == 1 {
		return nil
	}

	return ErrInvalidToken
}

// newToken creates a new Token with the given Config and gin.Context.
// It ensures that the session is valid before returning the token.
func newToken(cfg *Config, ctx *gin.Context) (*Token, error) {
	token := &Token{
		cfg: cfg,
		ctx: ctx,
	}
	err := token.ensureValidSession()

	return token, err
}

// Cookie generates an HTTP cookie from the token with specific security constraints.
func (r *Token) Cookie() *http.Cookie {
	// Set the cookie with the following constraints:
	// - `Domain` is not set, so the cookie is a host-only cookie.
	// - `Expires` is not set, so the cookie is a session cookie.
	// - `Secure` is set, so the cookie is only sent over HTTPS.
	// - `HttpOnly` is unset, so JavaScript can access the cookie.
	// - `MaxAge` is not set, so the cookie is a session cookie.
	// - `Partitioned` is set, so the cookie is sent in a separate cookie jar.
	return &http.Cookie{
		Name:        r.cfg.CookieName,
		Value:       url.QueryEscape(r.String()),
		Path:        "/",
		Secure:      true,
		HttpOnly:    false,
		SameSite:    r.cfg.CookieSameSite.SameSite,
		Partitioned: true,
	}
}

// HMAC generates a hash-based message authentication code for the token's
// session ID and random bytes using the configured secret and signer.
func (r *Token) HMAC() []byte {
	var msg bytes.Buffer

	msg.WriteString(strconv.Itoa(len(r.sessionID)))
	msg.WriteByte('!')
	msg.WriteString(r.sessionID)
	msg.WriteByte('!')
	msg.WriteString(strconv.Itoa(len(r.random)))
	msg.WriteByte('!')
	msg.Write(r.random)

	hash := hmac.New(r.cfg.Signer, r.cfg.Secret)
	hash.Write(msg.Bytes())

	return hash.Sum(nil)
}

// String returns the token as a base64 encoded string in the format "HMAC.Random".
func (r *Token) String() string {
	return base64.RawURLEncoding.EncodeToString(r.HMAC()) + "." +
		base64.RawURLEncoding.EncodeToString(r.random)
}

// ensureValidSession checks for a valid session in the context and validates its ID.
// It returns an error if no session is found or if the session ID is invalid.
func (r *Token) ensureValidSession() error {
	principal, ok := authentication.PrincipalFromContext(r.ctx)
	if !ok {
		return ErrNoSession
	}

	session, ok := principal.(authentication.SessionProvider)
	if !ok {
		return ErrNoSession
	}

	sessionID := session.SessionID()

	//nolint:mnd // OWASP recommends at least 128 bits of entropy for session identifiers.
	if len(sessionID) < 16 {
		return ErrInvalidSessionID
	}

	return nil
}

// parseTokenString splits a token string into its HMAC and random components, decoding them from base64 URL encoding.
// It returns the decoded HMAC and random byte slices along with any error encountered during the process.
func parseTokenString(tokenString string) (decodedHMAC, decodedRandom []byte, err error) {
	encodedHMAC, encodedRandom, ok := strings.Cut(tokenString, ".")
	if !ok {
		return nil, nil, ErrInvalidToken
	}

	decodedHMAC, err = base64.RawURLEncoding.DecodeString(encodedHMAC)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode HMAC: %w", err)
	}

	decodedRandom, err = base64.RawURLEncoding.DecodeString(encodedRandom)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode random bytes: %w", err)
	}

	return
}
