package authorization

import (
	"strings"
)

// Entitlement represents a type for entitlements, which are strings that encode both an action and a resource.
type Entitlement string

// Unpack splits the Entitlement into its action and resource parts.
// It returns the resource and action if the Entitlement is valid, otherwise it returns empty strings.
func (r Entitlement) Unpack() (string, Action) {
	//nolint:mnd // By definition, Entitlements always consist of two parts.
	parts := strings.SplitN(string(r), "_", 2)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		if _, ok := Actions[Action(parts[0])]; ok {
			return parts[1], Action(parts[0])
		}
	}
	return "", ""
}
