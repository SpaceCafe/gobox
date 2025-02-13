package saml

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gosaml/samlsp"
)

var (
	ErrNoSessionAttributes = errors.New("no session attributes found")
)

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
// an empty string if no such attributes is present.
func (r *Attributes) Get(name string) string {
	if key, ok := r.config.Mapping[name]; ok {
		return r.Attributes.Get(key)
	}
	return r.Attributes.Get(name)
}

// GetAll returns all mapped attributes named `name` or
// an empty []string if no such attributes is present.
func (r *Attributes) GetAll(name string) []string {
	if key, ok := r.config.Mapping[name]; ok {
		name = key
	}
	if r.Attributes == nil || len(r.Attributes[name]) == 0 {
		return []string{}
	}
	return r.Attributes[name]
}
