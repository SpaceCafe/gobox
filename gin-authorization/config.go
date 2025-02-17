package authorization

import (
	"errors"
	"path"
)

var (
	//nolint:gochecknoglobals // Used as default value and cannot be declared as constant due to its type.
	DefaultExcludedRoutes = make([]string, 0)
	//nolint:gochecknoglobals // Used as default value and cannot be declared as constant due to its type.
	DefaultRoles = make(map[string][]Entitlement)
	//nolint:gochecknoglobals // Used as default value and cannot be declared as constant due to its type.
	DefaultGroups = make(map[string][]Entitlement)
	//nolint:gochecknoglobals // Used as default value and cannot be declared as constant due to its type.
	DefaultUserAuthAdapter = NewJWTUserAuthAdapter

	ErrInvalidExcludedRoutes = errors.New("excluded routes must not be nil")
	ErrInvalidExcludedRoute  = errors.New("excluded route must be absolute and not end with a slash")
	ErrNoRoleMapper          = errors.New("role mapper must not be nil")
	ErrNoGroupMapper         = errors.New("group mapper must not be nil")
	ErrNoRoles               = errors.New("roles must not be nil")
	ErrNoGroups              = errors.New("groups must not be nil")
	ErrNoUserAuthAdapter     = errors.New("user auth adapter must not be nil")
)

// Config holds configuration related to user and API key authentication.
type Config struct {
	// ExcludedRoutes is a list of routes that are excluded from authorization checks.
	ExcludedRoutes []string `json:"excluded_routes" yaml:"excluded_routes" mapstructure:"excluded_routes"`

	// RoleMappingAdapter maps roles to entitlements.
	RoleMappingAdapter RoleGroupMappingAdapter

	// GroupMappingAdapter maps groups to entitlements.
	GroupMappingAdapter RoleGroupMappingAdapter

	// Roles is a map of role names to their entitlements.
	Roles map[string][]Entitlement `json:"roles" yaml:"roles" mapstructure:"roles"`

	// Groups is a map of group names to their entitlements.
	Groups map[string][]Entitlement `json:"groups" yaml:"groups" mapstructure:"groups"`

	// UserAuthAdapter is a function that returns a UserAuthAdapter for user authorization.
	UserAuthAdapter UserAuthAdapterFunc
}

// NewConfig creates and returns a new Config having default values.
func NewConfig() *Config {
	config := &Config{
		ExcludedRoutes:  DefaultExcludedRoutes,
		Roles:           DefaultRoles,
		Groups:          DefaultGroups,
		UserAuthAdapter: DefaultUserAuthAdapter,
	}
	config.RoleMappingAdapter = NewConfigRoleGroupMappingAdapter(func() map[string][]Entitlement { return config.Roles })
	config.GroupMappingAdapter = NewConfigRoleGroupMappingAdapter(func() map[string][]Entitlement { return config.Groups })

	return config
}

// Validate ensures the all necessary configurations are filled and within valid confines.
// Any misconfiguration results in well-defined standardized errors.
func (r *Config) Validate() error {
	if r.ExcludedRoutes == nil {
		return ErrInvalidExcludedRoutes
	}
	for i := range r.ExcludedRoutes {
		if !path.IsAbs(r.ExcludedRoutes[i]) {
			return ErrInvalidExcludedRoute
		}
	}
	if r.RoleMappingAdapter == nil {
		return ErrNoRoleMapper
	}
	if r.GroupMappingAdapter == nil {
		return ErrNoGroupMapper
	}
	if r.Roles == nil {
		return ErrNoRoles
	}
	if r.Groups == nil {
		return ErrNoGroups
	}
	if r.UserAuthAdapter == nil {
		return ErrNoUserAuthAdapter
	}
	return nil
}
