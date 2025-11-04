package authentication

// Principal defines an authenticated identity.
type Principal interface {
	ID() string
	Name() string
}

// MutablePrincipal defines methods for setting identity.
type MutablePrincipal interface {
	SetID(id string)
	SetName(name string)
}

// EmailProvider defines methods for retrieving email.
type EmailProvider interface {
	Email() string
}

// MutableEmailProvider defines methods for setting email.
type MutableEmailProvider interface {
	SetEmail(email string)
}

// PersonProvider defines methods for retrieving personal info.
type PersonProvider interface {
	FirstName() string
	LastName() string
}

// MutablePersonProvider defines methods for setting personal info.
type MutablePersonProvider interface {
	SetFirstName(firstName string)
	SetLastName(lastName string)
}

// AuthorizationProvider defines methods for retrieving authorization groups, roles, and entitlements.
type AuthorizationProvider interface {
	Groups() []string
	Roles() []string
	Entitlements() []string
}

// MutableAuthorizationProvider defines methods for setting authorization groups, roles, and entitlements.
type MutableAuthorizationProvider interface {
	SetGroups(groups []string)
	SetRoles(roles []string)
	SetEntitlements(entitlements []string)
}
