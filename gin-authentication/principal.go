package authentication

import (
	"time"
)

// Principal defines an authenticated identity.
type Principal interface {
	ID() string
	Name() string
}

// SessionProvider defines methods for retrieving session information.
type SessionProvider interface {
	Principal
	SessionID() string
	CreatedAt() time.Time
	ExpiresAt() time.Time
}

// MutablePrincipal defines methods for setting identity.
type MutablePrincipal interface {
	Principal
	SetID(id string)
	SetName(name string)
}

// EmailProvider defines methods for retrieving email.
type EmailProvider interface {
	Principal
	Email() string
}

// MutableEmailProvider defines methods for setting email.
type MutableEmailProvider interface {
	EmailProvider
	SetEmail(email string)
}

// PersonProvider defines methods for retrieving personal info.
type PersonProvider interface {
	Principal
	FirstName() string
	LastName() string
}

// MutablePersonProvider defines methods for setting personal info.
type MutablePersonProvider interface {
	PersonProvider
	SetFirstName(firstName string)
	SetLastName(lastName string)
}

// AuthorizationProvider defines methods for retrieving authorization groups, roles, and entitlements.
type AuthorizationProvider interface {
	Principal
	Groups() []string
	Roles() []string
	Entitlements() []string
}

// MutableAuthorizationProvider defines methods for setting authorization groups, roles, and entitlements.
type MutableAuthorizationProvider interface {
	AuthorizationProvider
	SetGroups(groups []string)
	SetRoles(roles []string)
	SetEntitlements(entitlements []string)
}
