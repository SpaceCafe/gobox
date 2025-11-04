package authentication

var (
	_ Principal = (*DefaultPrincipal)(nil)

	//nolint:gochecknoglobals // tokenPrincipal is the default principal used for API tokens.
	tokenPrincipal = &DefaultPrincipal{id: "token", name: "token"}
)

// DefaultPrincipal is the default implementation of the Principal interface.
type DefaultPrincipal struct {
	id   string
	name string
}

// ID returns the principal's id.
func (r *DefaultPrincipal) ID() string {
	return r.id
}

// Name returns the principal's name.
func (r *DefaultPrincipal) Name() string {
	return r.name
}
