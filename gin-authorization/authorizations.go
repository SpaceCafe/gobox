package authorization

import (
	"github.com/gin-gonic/gin"
)

// Authorizations is a map where the keys are resource names and the values are maps of action names to empty structs.
// Example:
//
//	Authorizations{"books": {"read": {}, "create": {}}, "cars": {{"update": {}}}
type Authorizations map[string]map[Action]struct{}

// NewAuthorizations initializes an Authorizations map based on user, group, and role authorizations.
// It merges authorizations from user groups and roles, and adds entitlements.
func NewAuthorizations(config *Config, ctx *gin.Context) Authorizations {

	// Retrieve initial authorizations from the user
	userAuth := config.UserAuthAdapter(ctx)
	authorizations := userAuth.Authorizations()

	// Merge group authorizations
	for _, group := range userAuth.Groups() {
		if groupAuth, ok := config.GroupMapper.Map()[group]; ok {
			authorizations.Merge(groupAuth)
		}
	}

	// Merge role authorizations
	for _, role := range userAuth.Roles() {
		if roleAuth, ok := config.RoleMapper.Map()[role]; ok {
			authorizations.Merge(roleAuth)
		}
	}

	// Add entitlements
	for _, entitlement := range userAuth.Entitlements() {
		authorizations.Add(Entitlement(entitlement).Unpack())
	}

	return authorizations
}

// Merge combines another Authorizations map into the current one.
// It adds actions from the other map to the existing resources or creates new resources if they don't exist.
func (r Authorizations) Merge(authorizations Authorizations) {
	for resource, actions := range authorizations {
		if existingActions, ok := r[resource]; ok {
			for action := range actions {
				existingActions[action] = struct{}{}
			}
		} else {
			r[resource] = actions
		}
	}
}

// Add inserts a new action for a resource into the Authorizations map.
// If the resource already exists, it adds the action to the existing actions.
func (r Authorizations) Add(resource string, action Action) {
	if existingActions, ok := r[resource]; ok {
		existingActions[action] = struct{}{}
	} else {
		r[resource] = map[Action]struct{}{action: {}}
	}
}

// IsAuthorized checks if a given action is authorized for a specified resource.
// It returns true if the action is allowed, otherwise false.
func (r Authorizations) IsAuthorized(resource string, action Action) bool {
	if existingActions, ok := r[resource]; ok {
		if _, ok := existingActions[action]; ok {
			return true
		}
	}

	// Check for global permissions
	if existingActions, ok := r["*"]; ok {
		if _, ok := existingActions[action]; ok {
			return true
		}
	}
	return false
}
