package authorization

// RoleGroupMap defines a mapping from role/group names to their respective authorizations.
type RoleGroupMap map[string]Authorizations

// RoleGroupMapper is an interface for mapping roles/groups to their authorizations.
type RoleGroupMapper interface {
	Map() RoleGroupMap
}

// ConfigRoleMapper maps roles to their authorizations using a configuration.
type ConfigRoleMapper struct {
	Config *Config
	cache  RoleGroupMap
}

// Map returns the cached role-to-authorization map or creates it if not already cached.
func (a *ConfigRoleMapper) Map() RoleGroupMap {
	if a.cache == nil {
		a.cache = createMap(a.Config.Roles)
	}
	return a.cache
}

// ConfigGroupMapper maps groups to their authorizations using a configuration.
type ConfigGroupMapper struct {
	Config *Config
	cache  RoleGroupMap
}

// Map returns the cached group-to-authorization map or creates it if not already cached.
func (a *ConfigGroupMapper) Map() RoleGroupMap {
	if a.cache == nil {
		a.cache = createMap(a.Config.Groups)
	}
	return a.cache
}

// createMap is a helper function to create the RoleGroupMap from a given map of strings to slices of Entitlement.
func createMap(source map[string][]Entitlement) RoleGroupMap {
	result := make(RoleGroupMap)
	for name, entitlements := range source {

		// Skip empty names or entitlements
		if name != "" && entitlements != nil {
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
