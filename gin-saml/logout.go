package saml

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/spacecafe/gosaml"
)

// RequestType represents different types of SAML logout requests and responses.
type RequestType int

const (
	LogoutRequestPost RequestType = 0 + iota
	LogoutRequestRedirect
	LogoutResponsePost
	LogoutResponseRedirect
)

var (
	ErrNoLogoutRequest = errors.New("no saml logout request found")
)

// validateLogoutRequest validates SAML logout requests and responses.
// It takes a gin context as input and returns the logout request or response data, type of request/response, and an error if any.
func (r *SAML) validateLogoutRequest(ctx *gin.Context) (string, RequestType, error) {
	// If POST request includes a SAML LogoutResponse:
	// validating the logout response using SP's function and returning relevant data or error.
	if data, ok := ctx.GetPostForm("SAMLResponse"); ok && len(data) > 0 {
		return data, LogoutResponsePost, r.middleware.ServiceProvider.ValidateLogoutResponseForm(data)
	}

	// If POST request includes a SAML LogoutRequest:
	// validating the logout request using SP's function and returning relevant data or error.
	if data, ok := ctx.GetPostForm("SAMLRequest"); ok && len(data) > 0 {
		return data, LogoutRequestPost, r.middleware.ServiceProvider.ValidateLogoutRequestForm(data)
	}

	// If GET request includes a SAML LogoutResponse:
	// validating the logout response using SP's function and returning relevant data or error.
	if data, ok := ctx.GetQuery("SAMLResponse"); ok && len(data) > 0 {
		return data, LogoutResponseRedirect, r.middleware.ServiceProvider.ValidateLogoutResponseRedirect(data)
	}

	// If GET request includes a SAML LogoutRequest:
	// validating the logout response using SP's function and returning relevant data or error.
	if data, ok := ctx.GetQuery("SAMLRequest"); ok && len(data) > 0 {
		return data, LogoutRequestRedirect, r.middleware.ServiceProvider.ValidateLogoutRequestRedirect(data)
	}

	// If no valid logout request is present in POST/GET data.
	return "", LogoutRequestPost, ErrNoLogoutRequest
}

// newLogoutRequest decodes the base64 encoded SAML logout request to a LogoutRequest struct.
func newLogoutRequest(data string) (*saml.LogoutRequest, error) {
	logoutRequest := new(saml.LogoutRequest)

	// Decoding the base64 encoded logout request.
	rawRequestBuf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// Parsing and decoding the XML format of logout request to LogoutRequest struct.
	err = xml.NewDecoder(bytes.NewBuffer(rawRequestBuf)).Decode(logoutRequest)
	return logoutRequest, err
}
