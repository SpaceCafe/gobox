package saml

import (
	"net/http"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
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
		r.config.Logger.Warn(err)
		_ = ctx.Error(problems.ProblemBadRequest)
		return
	}

	switch requestType {
	// If it's a logout response then redirect to configured post logout URI.
	case LogoutResponsePost, LogoutResponseRedirect:
		ctx.Redirect(http.StatusSeeOther, r.config.PostLogoutURI)

	// If it's a logout request then delete session and response to IdP.
	case LogoutRequestPost, LogoutRequestRedirect:
		// Delete the SAML session.
		if err := r.middleware.Session.DeleteSession(ctx.Writer, ctx.Request); err != nil {
			r.config.Logger.Warnf("unable to delete session: %v", err)
			_ = ctx.Error(problems.ProblemInternalError)
			return
		}

		// Create the logout request payload.
		logoutRequest, err := newLogoutRequest(data)
		if err != nil {
			r.config.Logger.Warnf("unable to parse logout request: %v", err)
			_ = ctx.Error(problems.ProblemBadRequest)
			return
		}

		// Create the response to the logout request.
		switch requestType {
		case LogoutRequestPost:
			response, err := r.middleware.ServiceProvider.MakePostLogoutResponse(logoutRequest.ID, "")
			if err != nil {
				r.config.Logger.Warnf("unable to build logout response: %v", err)
				_ = ctx.Error(problems.ProblemInternalError)
				return
			}

			// Write the HTML page with the logout response.
			ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			ctx.String(
				http.StatusOK,
				`<!DOCTYPE html><html><head><meta charset="utf-8" /></head><body><noscript><p><strong>Note:</strong> Since your browser does not support JavaScript, you must press the Continue button once to proceed.</p></noscript>%s</body></html>`,
				response,
			)

		case LogoutRequestRedirect:
		default:
			_ = ctx.Error(problems.ProblemInternalError)
		}

	default:
		_ = ctx.Error(problems.ProblemInternalError)
	}
}
