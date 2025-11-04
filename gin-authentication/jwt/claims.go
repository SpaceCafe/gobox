package jwt

import (
	jwt2 "github.com/golang-jwt/jwt/v5"
)

// IdentityClaims represents the standard OpenID Connect claims for user identity information.
// OpenID Connect Core 1.0, Section 5.1
// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
type IdentityClaims struct {
	// Name is the End-User's full name in displayable form including all name parts,
	// possibly including titles and suffixes, ordered according to the End-User's locale and preferences.
	Name string `json:"name,omitempty"`

	// GivenName is the given name(s) or first name(s) of the End-User.
	// Note that in some cultures, people can have multiple given names;
	// all can be present, with the names being separated by space characters.
	GivenName string `json:"given_name,omitempty"`

	// FamilyName is the surname(s) or last name(s) of the End-User.
	// Note that in some cultures, people can have multiple family names or no family name;
	// all can be present, with the names being separated by space characters.
	FamilyName string `json:"family_name,omitempty"`

	// MiddleName is the middle name(s) of the End-User.
	// Note that in some cultures, people can have multiple middle names;
	// all can be present, with the names being separated by space characters.
	// Also note that in some cultures, middle names are not used.
	MiddleName string `json:"middle_name,omitempty"`

	// Nickname is the casual name of the End-User that may or may not be the same as the given_name.
	// For instance, a nickname value of Mike might be returned alongside a given_name value of Michael.
	Nickname string `json:"nickname,omitempty"`

	// PreferredUsername is the shorthand name by which the End-User wishes to be referred to at the RP,
	// such as janedoe or j.doe. This value MAY be any valid JSON string including special characters
	// such as @, /, or whitespace. The RP MUST NOT rely upon this value being unique.
	PreferredUsername string `json:"preferred_username,omitempty"`

	// Profile is the URL of the End-User's profile page.
	// The contents of this Web page SHOULD be about the End-User.
	Profile string `json:"profile,omitempty"`

	// Picture is the URL of the End-User's profile picture.
	// This URL MUST refer to an image file (for example, a PNG, JPEG, or GIF image file),
	// rather than to a Web page containing an image.
	Picture string `json:"picture,omitempty"`

	// Website is the URL of the End-User's Web page or blog.
	// This Web page SHOULD contain information published by the End-User or
	// an organization that the End-User is affiliated with.
	Website string `json:"website,omitempty"`

	// Email is the End-User's preferred e-mail address.
	// Its value MUST conform to the RFC 5322 addr-spec syntax.
	// The RP MUST NOT rely upon this value being unique.
	Email string `json:"email,omitempty"`

	// Gender is the End-User's gender.
	// Values defined by this specification are female and male.
	// Other values MAY be used when neither of the defined values are applicable.
	Gender string `json:"gender,omitempty"`

	// Birthdate is the End-User's birthday, represented as an ISO 8601-1 YYYY-MM-DD format.
	// The year MAY be 0000, indicating that it is omitted.
	// To represent only the year, YYYY format is allowed.
	Birthdate string `json:"birthdate,omitempty"`

	// ZoneInfo is a string from IANA Time Zone Database representing the End-User's time zone.
	// For example, Europe/Paris or America/Los_Angeles.
	ZoneInfo string `json:"zoneinfo,omitempty"`

	// Locale is the End-User's locale, represented as a BCP47 language tag.
	// This is typically an ISO 639 Alpha-2 language code in lowercase and an ISO 3166-1 Alpha-2
	// country code in uppercase, separated by a dash. For example, en-US or fr-CA.
	Locale string `json:"locale,omitempty"`

	// PhoneNumber is the End-User's preferred telephone number.
	// E.164 is RECOMMENDED as the format of this Claim, for example, +1 (425) 555-1212 or +56 (2) 687 2400.
	PhoneNumber string `json:"phone_number,omitempty"`

	// Address is the End-User's preferred postal address.
	// The value of the address member is a JSON structure containing some or all of the members
	// defined in Section 5.1.1.
	Address *AddressClaim `json:"address,omitempty"`

	// UpdatedAt is the time the End-User's information was last updated.
	// Its value is a JSON number representing the number of seconds from 1970-01-01T00:00:00Z
	// as measured in UTC until the date/time.
	UpdatedAt int64 `json:"updated_at,omitempty"`

	// EmailVerified is true if the End-User's e-mail address has been verified; otherwise false.
	// When this Claim Value is true, this means that the OP took affirmative steps to ensure
	// that this e-mail address was controlled by the End-User at the time the verification was performed.
	EmailVerified bool `json:"email_verified,omitempty"`

	// PhoneNumberVerified is true if the End-User's phone number has been verified; otherwise false.
	// When this Claim Value is true, this means that the OP took affirmative steps to ensure
	// that this phone number was controlled by the End-User at the time the verification was performed.
	PhoneNumberVerified bool `json:"phone_number_verified,omitempty"`
}

// AddressClaim represents the End-User's preferred postal address.
// OpenID Connect Core 1.0, Section 5.1.1
// https://openid.net/specs/openid-connect-core-1_0.html#AddressClaim
type AddressClaim struct {
	// Formatted is the full mailing address, formatted for display or use on a mailing label.
	// This field MAY contain multiple lines, separated by newlines.
	Formatted string `json:"formatted,omitempty"`

	// StreetAddress is the full street address component, which MAY include house number,
	// street name, Post Office Box, and multi-line extended street address information.
	// This field MAY contain multiple lines, separated by newlines.
	StreetAddress string `json:"street_address,omitempty"`

	// Locality is the city or locality component.
	Locality string `json:"locality,omitempty"`

	// Region is the state, province, prefecture, or region component.
	Region string `json:"region,omitempty"`

	// PostalCode is the zip code or postal code component.
	PostalCode string `json:"postal_code,omitempty"`

	// Country is the country name component.
	Country string `json:"country,omitempty"`
}

// AuthorizationClaims represents the authorization-related claims for access tokens.
// RFC 9068 - JSON Web Token (JWT) Profile for OAuth 2.0 Access Tokens
// https://www.rfc-editor.org/rfc/rfc9068.html
type AuthorizationClaims struct {
	// Roles contain user or client associated roles.
	// RFC 9068 Section 2.2
	// https://www.rfc-editor.org/rfc/rfc9068.html#section-2.2.3.1
	Roles []string `json:"roles,omitempty"`

	// Groups contain user or client associated groups.
	// RFC 9068 Section 2.2
	// https://www.rfc-editor.org/rfc/rfc9068.html#section-2.2.3.1
	Groups []string `json:"groups,omitempty"`

	// Entitlements contain user or client associated entitlements.
	// RFC 9068 Section 2.2
	// https://www.rfc-editor.org/rfc/rfc9068.html#section-2.2
	Entitlements []string `json:"entitlements,omitempty"`

	// AuthorizationDetails contain detailed authorization information as defined in RFC 9396.
	// This enables fine-grained authorization by expressing the specifics about the access
	// being requested or granted.
	// RFC 9396 Section 9.1
	// https://www.rfc-editor.org/rfc/rfc9396.html#section-9.1
	AuthorizationDetails []AuthorizationDetail `json:"authorization_details,omitempty"`
}

// AuthorizationDetail represents detailed authorization information as defined in RFC 9396.
// RFC 9396 Section 9.1
// https://www.rfc-editor.org/rfc/rfc9396.html#section-9.1
type AuthorizationDetail struct {
	Type    string   `json:"type"`
	Actions []string `json:"actions"`
}

// Claims represents the claims for a JWT token.
type Claims struct {
	IdentityClaims
	jwt2.RegisteredClaims
	AuthorizationClaims
}

func NewClaims(subject string) *Claims {
	return &Claims{
		RegisteredClaims: jwt2.RegisteredClaims{
			Subject: subject,
		},
	}
}

func (r *Claims) Email() string {
	return r.IdentityClaims.Email
}

func (r *Claims) Entitlements() []string {
	return r.AuthorizationClaims.Entitlements
}

func (r *Claims) FirstName() string {
	return r.GivenName
}

func (r *Claims) Groups() []string {
	return r.AuthorizationClaims.Groups
}

func (r *Claims) ID() string {
	return r.Subject
}

func (r *Claims) LastName() string {
	return r.FamilyName
}

func (r *Claims) Name() string {
	return r.IdentityClaims.Name
}

func (r *Claims) Roles() []string {
	return r.AuthorizationClaims.Roles
}

func (r *Claims) SetEmail(email string) {
	r.IdentityClaims.Email = email
}

func (r *Claims) SetEntitlements(entitlements []string) {
	r.AuthorizationClaims.Entitlements = entitlements
}

func (r *Claims) SetFirstName(firstName string) {
	r.GivenName = firstName
}

func (r *Claims) SetGroups(groups []string) {
	r.AuthorizationClaims.Groups = groups
}

func (r *Claims) SetID(id string) {
	r.Subject = id
}

func (r *Claims) SetLastName(lastName string) {
	r.FamilyName = lastName
}

func (r *Claims) SetName(name string) {
	r.IdentityClaims.Name = name
}

func (r *Claims) SetRoles(roles []string) {
	r.AuthorizationClaims.Roles = roles
}
