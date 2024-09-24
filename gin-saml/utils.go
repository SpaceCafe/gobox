package saml

import (
	"errors"

	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
)

var (
	ErrNoSessionAttributes = errors.New("no session attributes found")
)

type Attributes struct {
	samlsp.Attributes
	config *Config
}

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

// Get returns the first mapped attribute named `name` or an empty string if
// no such attributes is present.
func (r *Attributes) Get(name string) string {
	if key, ok := r.config.Mapping[name]; ok {
		return r.Attributes.Get(key)
	}
	return r.Attributes.Get(name)
}
