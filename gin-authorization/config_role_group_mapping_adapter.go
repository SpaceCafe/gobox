package authorization

import (
	"strings"
)

// ConfigRoleGroupMappingAdapter maps roles and groups to their authorizations using a configuration.
type ConfigRoleGroupMappingAdapter struct {
	config func() map[string][]Entitlement
	cache  RoleGroupMap
}

// NewConfigRoleGroupMappingAdapter initializes a new instance of ConfigRoleGroupMappingAdapter using the provided configuration function.
// A factory function is necessary as the role and group mappings are not initialized at application startup.
func NewConfigRoleGroupMappingAdapter(config func() map[string][]Entitlement) *ConfigRoleGroupMappingAdapter {
	return &ConfigRoleGroupMappingAdapter{
		config: config,
	}
}

// Map returns the cached role-to-authorization map or creates it if not already cached.
func (a *ConfigRoleGroupMappingAdapter) Map() RoleGroupMap {
	if a.cache == nil {
		a.cache = createMap(a.config())
	}
	return a.cache
}

// createMap is a helper function to create the RoleGroupMap from a given map of strings to slices of Entitlement.
func createMap(source map[string][]Entitlement) RoleGroupMap {
	result := make(RoleGroupMap)
	for name, entitlements := range source {

		// Skip empty names or entitlements
		if name != "" && entitlements != nil {
			name = strings.ToLower(name)
			auth := make(Authorizations)
			for _, entitlement := range entitlements {

				// Parse the entitlement to get the action and resource.
				if resource, action := entitlement.Unpack(); action != "" && resource != "" {
					if _, ok := auth[resource]; !ok {
						auth[resource] = make(map[Action]struct{})
					}
					auth[resource][action] = struct{}{}
				}
			}

			// Only add non-empty authorizations.
			if len(auth) > 0 {
				result[name] = auth
			}
		}
	}
	return result
}
