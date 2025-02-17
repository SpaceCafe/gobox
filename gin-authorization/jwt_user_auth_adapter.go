package authorization

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/spacecafe/gobox/gin-jwt"
)

// JWTUserAuthAdapter is a struct that adapts JWT claims to user authentication details.
type JWTUserAuthAdapter struct {
	claims jwt.AuthorizationClaims
}

// NewJWTUserAuthAdapter creates a new instance of JWTUserAuthAdapter by extracting claims from the given context.
func NewJWTUserAuthAdapter(ctx *gin.Context) UserAuthAdapter {
	return &JWTUserAuthAdapter{
		claims: jwt.GetClaims(ctx),
	}
}

// Roles returns the roles associated with the user based on the JWT claims.
func (r *JWTUserAuthAdapter) Roles() []string {
	if r.claims != nil {
		return r.claims.GetRoles()
	}
	return []string{}
}

// Groups returns the groups associated with the user based on the JWT claims.
func (r *JWTUserAuthAdapter) Groups() []string {
	if r.claims != nil {
		return r.claims.GetGroups()
	}
	return []string{}
}

// Entitlements returns the entitlements associated with the user based on the JWT claims.
func (r *JWTUserAuthAdapter) Entitlements() []string {
	if r.claims != nil {
		return r.claims.GetEntitlements()
	}
	return []string{}
}

// Authorizations constructs a map of authorizations from the JWT claims, filtering by known actions.
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
