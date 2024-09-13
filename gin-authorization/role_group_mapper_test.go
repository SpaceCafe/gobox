package authorization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigRoleGroupMapper_Map(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   RoleGroupMap
	}{
		{
			"valid group map",
			&Config{
				Roles:  map[string][]Entitlement{"customer": {"read_books", "list_books"}},
				Groups: map[string][]Entitlement{"customer": {"read_books", "list_books"}},
			},
			RoleGroupMap{"customer": {"books": {ReadAction: {}, ListAction: {}}}},
		},
		{
			"invalid action",
			&Config{
				Roles:  map[string][]Entitlement{"customer": {"invalid_books", "list_books"}},
				Groups: map[string][]Entitlement{"customer": {"invalid_books", "list_books"}},
			},
			RoleGroupMap{"customer": {"books": {ListAction: {}}}},
		},
		{
			"invalid resource",
			&Config{
				Roles:  map[string][]Entitlement{"customer": {"read"}},
				Groups: map[string][]Entitlement{"customer": {"read"}},
			},
			RoleGroupMap{},
		},
		{
			"invalid group name",
			&Config{
				Roles:  map[string][]Entitlement{"": {"read_books"}},
				Groups: map[string][]Entitlement{"": {"read_books"}},
			},
			RoleGroupMap{},
		},
		{
			"empty group",
			&Config{
				Roles:  map[string][]Entitlement{"customer": {}},
				Groups: map[string][]Entitlement{"customer": {}},
			},
			RoleGroupMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoles := &ConfigRoleMapper{Config: tt.config}
			gotGroups := &ConfigGroupMapper{Config: tt.config}
			assert.Equal(t, tt.want, gotRoles.Map())
			assert.Equal(t, tt.want, gotGroups.Map())
		})
	}
}
