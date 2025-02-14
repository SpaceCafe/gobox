package saml

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gosaml/samlsp"
)

var (
	ErrNoSessionAttributes = errors.New("no session attributes found")
)

type ErrNoAttribute struct {
	Name string
}

func (e ErrNoAttribute) Error() string {
	return "no attribute '" + e.Name + "' found"
}

// Attributes is a type that embeds samlsp.Attributes and includes a configuration pointer.
type Attributes struct {
	samlsp.Attributes
	config *Config
}

// GetAttributes retrieves SAML attributes from the session context.
// It returns an error if no session attributes are found.
func (r *SAML) GetAttributes(ctx *gin.Context) (attributes *Attributes, err error) {
	session, _ := ctx.Get("saml.session")
	sessionAttributes, ok := session.(samlsp.SessionWithAttributes)
	if !ok {
		err = ErrNoSessionAttributes
		return
	}
	attributes = &Attributes{sessionAttributes.GetAttributes(), r.config}
	return
}

// Get returns the first mapped attribute named `name` or
// an empty string if no such attribute is present.
func (r *Attributes) Get(name string) (value string) {
	value, _ = r.MustGet(name)
	return
}

// GetAll returns all mapped attributes named `name` or
// an empty []string if no such attributes is present.
func (r *Attributes) GetAll(name string) (values []string) {
	values, _ = r.MustGetAll(name)
	return
}

// MustGet returns the first mapped attribute named `name` or
// an error if no such attribute is present.
func (r *Attributes) MustGet(name string) (value string, err error) {
	if key, ok := r.config.Mapping[name]; ok {
		name = key
	}
	value = r.Attributes.Get(name)
	if value == "" {
		err = ErrNoAttribute{Name: name}
	}
	return
}

// MustGetAll returns all mapped attributes named `name` or
// an error if no such attribute is present.
func (r *Attributes) MustGetAll(name string) (values []string, err error) {
	if key, ok := r.config.Mapping[name]; ok {
		name = key
	}
	if r.Attributes == nil || len(r.Attributes[name]) == 0 {
		return []string{}, ErrNoAttribute{Name: name}
	}
	return r.Attributes[name], nil
}
