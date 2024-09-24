package saml

import (
	"errors"

	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
)

var (
	ErrNoSessionAttributes = errors.New("no session attributes found")
)

func GetAttributes(ctx *gin.Context) (attributes samlsp.Attributes, err error) {
	session, _ := ctx.Get("saml.session")
	sessionAttributes, ok := session.(samlsp.SessionWithAttributes)
	if !ok {
		err = ErrNoSessionAttributes
		return
	}
	attributes = sessionAttributes.GetAttributes()
	return
}
