package authorization

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/spacecafe/gobox/gin-jwt"
)

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

type JWTUserAuthAdapter struct {
	claims jwt.AuthorizationClaims
}

func NewJWTUserAuthAdapter(ctx *gin.Context) UserAuthAdapter {
	result := &JWTUserAuthAdapter{}
	if token, ok := ctx.Get("jwt/token"); ok {
		result.claims = token.(*jwt.Token).Claims.(jwt.AuthorizationClaims)
	}
	return result
}

func (r *JWTUserAuthAdapter) Roles() []string {
	if r.claims != nil {
		return r.claims.GetRoles()
	}
	return []string{}
}

func (r *JWTUserAuthAdapter) Groups() []string {
	if r.claims != nil {
		return r.claims.GetGroups()
	}
	return []string{}
}

func (r *JWTUserAuthAdapter) Entitlements() []string {
	if r.claims != nil {
		return r.claims.GetEntitlements()
	}
	return []string{}
}

func (r *JWTUserAuthAdapter) Authorizations() Authorizations {
	result := make(Authorizations)
	if r.claims != nil {
		for _, detail := range r.claims.GetAuthorizationDetails() {
			if detail.Type != "" {
				for _, action := range detail.Actions {
					if _, ok := Actions[Action(action)]; ok {
						if _, ok := result[detail.Type]; !ok {
							result[detail.Type] = make(map[Action]struct{})
						}
						result[detail.Type][Action(action)] = struct{}{}
					}
				}
			}
		}
	}
	return result
}
