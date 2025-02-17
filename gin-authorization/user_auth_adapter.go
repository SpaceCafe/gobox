package authorization

import (
	"github.com/gin-gonic/gin"
)

// UserAuthAdapterFunc represents a function that accepts a Gin context and yields a UserAuthAdapter.
// This factory method produces a fresh instance of UserAuthAdapter for every incoming request.
type UserAuthAdapterFunc func(*gin.Context) UserAuthAdapter

type UserAuthAdapter interface {
	// Groups returns a slice of group names the user belongs to.
	// Example: ["group1", "group2"]
	Groups() []string

	// Roles returns a slice of role names assigned to the user
	// Example: ["admin", "guest"]
	Roles() []string

	// Entitlements returns a slice of user's entitlements in the format "<action>_<resource>".
	// Example: ["read_resource1", "write_resource2"]
	Entitlements() []string

	// Authorizations returns a map of user's authorizations with the schema: "<resource>[<action>]".
	// The outer map's keys are resource names, and the values are maps where:
	// - The keys are action names.
	// - The values are empty structs.
	// Example: {"resource1": {"read": {}, "write": {}}, "resource2": {"read": {}}}
	Authorizations() Authorizations
}
