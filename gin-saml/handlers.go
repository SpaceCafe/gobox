package saml

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	saml "github.com/spacecafe/gosaml"
)

const (
	LogoutRequestHTMLHead = `<html lang="en"><head><meta charset="utf-8"/><title>SAML Logout</title></head><body><noscript><p><strong>Note:</strong> Since your browser does not support JavaScript, you must press the Continue button once to proceed.</p></noscript>`
	LogoutRequestHTMLFoot = `</body></html>`
)

var (
	//nolint:gochecknoglobals // Used as http response value and cannot be declared as constant due to its type.
	logoutRequestHTMLHead = []byte(LogoutRequestHTMLHead)
	//nolint:gochecknoglobals // Used as http response value and cannot be declared as constant due to its type.
	logoutRequestHTMLFoot = []byte(LogoutRequestHTMLFoot)
)

// getMetadata retrieves SAML metadata using middleware's ServeMetadata method.
func (r *SAML) getMetadata(ctx *gin.Context) {
	r.middleware.ServeMetadata(ctx.Writer, ctx.Request)
}

// handleACS handles Assertion Consumer GetService (ACS) using middleware's ServeACS method.
func (r *SAML) handleACS(ctx *gin.Context) {
	r.middleware.ServeACS(ctx.Writer, ctx.Request)
}

// handleSLO handles Single Log Out (SLO).
func (r *SAML) handleSLO(ctx *gin.Context) {
	// Validate the request and get its type (logout response or logout request).
	// Parsing the form data / query param from HTTP Request.
	data, requestType, err := r.validateLogoutRequest(ctx)
	if err != nil {
		_ = ctx.Error(err)
		problems.ProblemBadRequest.Abort(ctx)
		return
	}

	// Delete the SAML session.
	if err = r.DeleteSession(ctx); err != nil {
		_ = ctx.Error(err)
		return
	}

	switch requestType {
	// If it's a logout response then redirect to configured post logout URI.
	case LogoutResponsePost, LogoutResponseRedirect:
		ctx.Redirect(http.StatusSeeOther, r.config.PostLogoutURI)

	// If it's a logout request then response to IdP.
	case LogoutRequestPost, LogoutRequestRedirect:
		// Create the logout request payload.
		var logoutRequest *saml.LogoutRequest
		if logoutRequest, err = newLogoutRequest(data); err != nil {
			_ = ctx.Error(err)
			problems.ProblemBadRequest.Abort(ctx)
			return
		}

		if requestType == LogoutRequestPost {
			// Create the response to the POST logout request.
			var response []byte
			response, err = r.middleware.ServiceProvider.MakePostLogoutResponse(logoutRequest.ID, "")
			if err != nil {
				_ = ctx.Error(err)
				problems.ProblemInternalError.Abort(ctx)
				return
			}

			// Write the HTML page with the logout response.
			ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			ctx.Status(http.StatusOK)
			_, _ = ctx.Writer.Write(logoutRequestHTMLHead)
			_, _ = ctx.Writer.Write(response)
			_, _ = ctx.Writer.Write(logoutRequestHTMLFoot)
		} else {
			// Create the redirection for the logout request.
			var redirectURL *url.URL
			redirectURL, err = r.middleware.ServiceProvider.MakeRedirectLogoutResponse(logoutRequest.ID, "")
			if err != nil {
				_ = ctx.Error(err)
				problems.ProblemInternalError.Abort(ctx)
				return
			}
			ctx.Redirect(http.StatusSeeOther, redirectURL.String())
		}
	default:
		_ = ctx.Error(problems.ProblemInternalError)
	}
}
