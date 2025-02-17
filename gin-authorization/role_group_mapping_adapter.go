package authorization

// RoleGroupMap defines a mapping from role or group names to their respective authorizations.
type RoleGroupMap map[string]Authorizations

// RoleGroupMappingAdapter is an interface for mapping roles or groups to their authorizations.
type RoleGroupMappingAdapter interface {
	Map() RoleGroupMap
}
