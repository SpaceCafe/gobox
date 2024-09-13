package authorization

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ResourceAction represents an action performed on a resource.
type ResourceAction struct {

	// Resource specifies the name or identifier of the resource.
	Resource string

	// Action specifies the action to be performed on the resource.
	Action Action
}

// NewResourceActions determines the actions to be taken on resources based on the HTTP method and URL path.
// For a GET request to "/users/:id/posts", it might return:
//
//	[]ResourceAction{{Resource: "posts", Action: "list"}, {Resource: "users", Action: "read"}}
func NewResourceActions(ctx *gin.Context, fullPath string) []ResourceAction {
	var result []ResourceAction

	segments := strings.Split(fullPath, "/")

	// Determine the primary and parent actions based on the HTTP method
	action, parentAction := getActions(ctx.Request.Method)

	// Loop in reverse order to build the result slice
	for i := len(segments) - 1; i >= 0; i-- {
		segment := segments[i]

		if strings.HasPrefix(segment, ":") {
			if len(result) == 0 && action == ListAction {
				action = ReadAction
			} else if parentAction == ListAction {
				parentAction = ReadAction
			}
			continue
		}

		// Append the appropriate ResourceAction to the result slice
		if len(result) == 0 {
			result = []ResourceAction{{Resource: segment, Action: action}}
		} else {
			result = append(result, ResourceAction{Resource: segment, Action: parentAction})
		}
	}

	return result
}

// getActions returns the primary and parent actions based on the HTTP method.
// The primary action is the action to be performed on the resource itself.
// The parent action is the action to be performed on the parent resource if a subresource is present.
func getActions(method string) (primaryAction, parentAction Action) {
	switch method {
	case "GET", "HEAD":
		return ListAction, ListAction
	case "POST":
		return CreateAction, UpdateAction
	case "PUT", "PATCH":
		return UpdateAction, UpdateAction
	case "DELETE":
		return DeleteAction, UpdateAction
	case "OPTIONS":
		return CapabilitiesAction, CapabilitiesAction
	default:
		return UnknownAction, UnknownAction
	}
}
