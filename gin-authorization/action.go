package authorization

// Action represents the type of action that can be performed on a resource.
type Action string

const (
	// CapabilitiesAction represents the action of retrieving supported operations on a resource.
	CapabilitiesAction Action = "capabilities"

	// CreateAction represents the action of creating a new resource.
	CreateAction Action = "create"

	// DeleteAction represents the action of deleting an existing resource.
	DeleteAction Action = "delete"

	// ReadAction represents the action of reading or retrieving a single resource.
	ReadAction Action = "read"

	// ListAction represents the action of listing multiple resources.
	ListAction Action = "list"

	// UnknownAction represents an undefined or unsupported action.
	UnknownAction Action = "unknown"

	// UpdateAction represents the action of updating an existing resource.
	UpdateAction Action = "update"
)

var (
	// Actions is a map that holds a set of predefined actions.
	// This map is used to check if a given action is valid.
	//nolint:gochecknoglobals // Maintain a set of predefined actions that are used throughout the application.
	Actions = map[Action]struct{}{
		CapabilitiesAction: {},
		CreateAction:       {},
		DeleteAction:       {},
		ReadAction:         {},
		ListAction:         {},
		UnknownAction:      {},
		UpdateAction:       {},
	}
)
