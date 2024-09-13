package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntitlement_Map(t *testing.T) {
	tests := []struct {
		name         string
		arg          Entitlement
		wantAction   Action
		wantResource string
	}{
		{
			name:         "create action",
			arg:          "create_example",
			wantAction:   CreateAction,
			wantResource: "example",
		},
		{
			name:         "delete action",
			arg:          "delete_example",
			wantAction:   DeleteAction,
			wantResource: "example",
		},
		{
			name:         "read action",
			arg:          "read_example",
			wantAction:   ReadAction,
			wantResource: "example",
		},
		{
			name:         "list action",
			arg:          "list_example",
			wantAction:   ListAction,
			wantResource: "example",
		},
		{
			name:         "update action",
			arg:          "update_example",
			wantAction:   UpdateAction,
			wantResource: "example",
		},
		{
			name:         "empty entitlement",
			arg:          "",
			wantAction:   "",
			wantResource: "",
		},
		{
			name:         "empty action",
			arg:          "_example",
			wantAction:   "",
			wantResource: "",
		},
		{
			name:         "empty resource",
			arg:          "create_",
			wantAction:   "",
			wantResource: "",
		},
		{
			name:         "invalid action",
			arg:          "invalid_example",
			wantAction:   "",
			wantResource: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResource, gotAction := tt.arg.Unpack()
			assert.Equal(t, tt.wantAction, gotAction)
			assert.Equal(t, tt.wantResource, gotResource)
		})
	}
}
