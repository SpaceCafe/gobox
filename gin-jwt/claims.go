package jwt

import (
	"time"

	jwt_ "github.com/golang-jwt/jwt/v5"
)

// IdentityClaims represents the interface for identity-related claims.
// It provides methods to retrieve the user's name, given name, and family name.
type IdentityClaims interface {

	// GetName returns the user's full name.
	GetName() (string, error)

	// GetGivenName returns the user's given name.
	GetGivenName() (string, error)

	// GetFamilyName returns the user's family name.
	GetFamilyName() (string, error)
}

// AuthorizationClaims represents the interface for authorization-related claims.
// It provides methods to retrieve the user's email, roles, groups, entitlements, and detailed authorization information.
type AuthorizationClaims interface {

	// GetEmail returns the user's primary email address.
	GetEmail() (string, error)

	// GetRoles returns the user's roles.
	GetRoles() []string

	// GetGroups returns the user's groups.
	GetGroups() []string

	// GetEntitlements returns the user's entitlements.
	GetEntitlements() []string

	// GetAuthorizationDetails returns the user's detailed authorization information.
	GetAuthorizationDetails() []AuthorizationDetail
}

// Claims represents the JWT claims for a user.
// It includes standard registered claims, as well as custom claims for roles, groups, entitlements, and authorization details.
type Claims struct {
	jwt_.RegisteredClaims

	// Roles contain user associated roles.
	// RFC 9068 Section 2.2
	// https://www.rfc-editor.org/rfc/rfc9068.html#section-2.2
	Roles []string `json:"roles"`

	// Groups contain user associated groups.
	// RFC 9068 Section 2.2
	// https://www.rfc-editor.org/rfc/rfc9068.html#section-2.2
	Groups []string `json:"groups"`

	// Entitlements contain user associated entitlements.
	// RFC 9068 Section 2.2
	// https://www.rfc-editor.org/rfc/rfc9068.html#section-2.2
	Entitlements []string `json:"entitlements"`

	// AuthorizationDetails contain detailed authorization information.
	// RFC 9396 Section 9.1
	// https://www.rfc-editor.org/rfc/rfc9396.html#section-9.1
	AuthorizationDetails []AuthorizationDetail `json:"authorization_details"`

	// Name is the user's full name.
	// OpenID Connect Core 1.0, Section 5.1
	// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
	Name string `json:"name"`

	// GivenName is the user's given name.
	// OpenID Connect Core 1.0, Section 5.1
	// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
	GivenName string `json:"given_name"`

	// FamilyName is the user's family name.
	// OpenID Connect Core 1.0, Section 5.1
	// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
	FamilyName string `json:"family_name"`

	// Email is the user's primary email address.
	// OpenID Connect Core 1.0, Section 5.1
	// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
	Email string `json:"email"`
}

// AuthorizationDetail represents detailed authorization information for a user.
// It includes the type of authorization and the actions allowed.
type AuthorizationDetail struct {
	Type    string   `json:"type"`
	Actions []string `json:"actions"`
}

// NewClaims creates a new Claims object with the provided configuration and subject.
// It initializes the RegisteredClaims with the issuer, subject, audience, expiration, not before, and issued at times.
func NewClaims(config *Config, subject string) *Claims {
	return &Claims{
		RegisteredClaims: jwt_.RegisteredClaims{
			Issuer:    config.Issuer,
			Subject:   subject,
			Audience:  config.Audiences,
			ExpiresAt: jwt_.NewNumericDate(time.Now().Add(config.TokenExpiration)),
			NotBefore: jwt_.NewNumericDate(time.Now()),
			IssuedAt:  jwt_.NewNumericDate(time.Now()),
		},
	}
}

// GetName returns the user's full name.
func (r *Claims) GetName() (string, error) {
	return r.Name, nil
}

// GetGivenName returns the user's given name.
func (r *Claims) GetGivenName() (string, error) {
	return r.GivenName, nil
}

// GetFamilyName returns the user's family name.
func (r *Claims) GetFamilyName() (string, error) {
	return r.FamilyName, nil
}

// GetEmail returns the user's primary email address.
func (r *Claims) GetEmail() (string, error) {
	return r.Email, nil
}

// GetRoles returns the user's roles.
func (r *Claims) GetRoles() []string {
	return r.Roles
}

// GetGroups returns the user's groups.
func (r *Claims) GetGroups() []string {
	return r.Groups
}

// GetEntitlements returns the user's entitlements.
func (r *Claims) GetEntitlements() []string {
	return r.Entitlements
}

// GetAuthorizationDetails returns the user's authorization details.
func (r *Claims) GetAuthorizationDetails() []AuthorizationDetail {
	return r.AuthorizationDetails
}
