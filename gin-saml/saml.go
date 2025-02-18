package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	saml "github.com/spacecafe/gosaml"
	"github.com/spacecafe/gosaml/samlsp"
)

const (
	FetchMetadataTimeout     = time.Minute
	ErrInvalidCertificates   = "certificates are not valid"
	ErrInvalidMetadataURL    = "metadata url is not valid"
	ErrFetchMetadata         = "could not fetch metadata"
	ErrInvalidApplicationURL = "application's url is not valid"
)

type SAML struct {
	config     *Config
	middleware *samlsp.Middleware
}

func New(config *Config, rg *gin.RouterGroup) (*SAML, error) {
	s := &SAML{
		config: config,
	}
	err := s.newMiddleware(rg.BasePath())
	if err != nil {
		return nil, err
	}

	// Add routes to gin.
	rg.GET("/saml/metadata", s.getMetadata)
	rg.Match([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch}, "/saml/acs", s.handleACS)
	rg.Match([]string{http.MethodGet, http.MethodPost}, "/saml/slo", s.handleSLO)

	return s, nil
}

func (r *SAML) newMiddleware(basePath string) error {
	// Allow login requests to be valid for a longer period of time.
	// Otherwise, users will receive errors during the authentication flow because the tracked request has expired.
	saml.MaxIssueDelay = r.config.MaxIssueDelay

	r.config.Logger.Debugf("load X509 certificate key pair: '%s' + '%s", r.config.CertFile, r.config.KeyFile)
	keyPair, err := tls.LoadX509KeyPair(r.config.CertFile, r.config.KeyFile)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidCertificates, err)
	}

	r.config.Logger.Debugf("parse X509 certificate leaf")
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidCertificates, err)
	}

	// Parse the IDP Metadata URL.
	idpMetadataURL, err := url.Parse(r.config.IDPMetadataURL)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidMetadataURL, err)
	}

	r.config.Logger.Debugf("fetch IDP metadata from '%s'", r.config.IDPMetadataURL)
	ctx, cancel := context.WithTimeout(context.Background(), FetchMetadataTimeout)
	defer cancel()

	idpMetadata, err := samlsp.FetchMetadata(ctx, http.DefaultClient, *idpMetadataURL)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrFetchMetadata, err)
	}
	r.config.Logger.Debugf("current metadata: %+v", idpMetadata)

	// Trailing slash is necessary to ensure correct URL for metadata, SLO, and ACS endpoints.
	if basePath != "" && basePath[len(basePath)-1] != '/' {
		basePath += "/"
	}

	// Parse the base path to set URL in SAML.
	baseURL, err := url.Parse(r.config.URI + basePath)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidApplicationURL, err)
	}

	// Create a new SAML GetService Provider.
	middleware, _ := samlsp.New(samlsp.Options{
		EntityID:            r.config.EntityID,
		URL:                 *baseURL,
		Key:                 keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:         keyPair.Leaf,
		AllowIDPInitiated:   r.config.AllowIDPInitiated,
		DefaultRedirectURI:  r.config.DefaultRedirectURI,
		IDPMetadata:         idpMetadata,
		SignRequest:         r.config.SignRequest,
		UseArtifactResponse: r.config.UseArtifactResponse,
		ForceAuthn:          r.config.ForceAuthn,
		CookieSameSite:      parseCookieSameSite(r.config.CookieSameSite),
		CookieName:          r.config.CookieName,
		LogoutBindings:      r.config.LogoutBindings,
	})
	middleware.ServiceProvider.AuthnNameIDFormat = saml.NameIDFormat(r.config.AuthnNameIDFormat)
	r.middleware = middleware
	r.config.Logger.Debugf("metadata url: %s", middleware.ServiceProvider.MetadataURL.String())

	return nil
}

// SetOnError changes the OnError handler for middleware.
func (r *SAML) SetOnError(function func(http.ResponseWriter, *http.Request, error)) {
	r.middleware.OnError = function
}

// RequireAccount is HTTP middleware that requires that each request be associated with a valid session.
// If the request is not associated with a valid session, then rather than serve the request,
// the middleware redirects the user to start the SAML auth flow.
func (r *SAML) RequireAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := r.middleware.Session.GetSession(ctx.Request)
		if session != nil {
			ctx.Set("saml.session", session)
			ctx.Next()
			return
		}
		if errors.Is(err, samlsp.ErrNoSession) {
			r.middleware.HandleStartAuthFlow(ctx.Writer, ctx.Request)
			ctx.Abort()
			return
		}

		r.middleware.OnError(ctx.Writer, ctx.Request, err)
		ctx.Abort()
	}
}

// onError handles errors during processing of SAML requests/responses.
func (r *SAML) onError(writer http.ResponseWriter, request *http.Request, err error) {
	r.config.Logger.Warnf("could not process saml request/response: %s", err)

	var invalidResponseError *saml.InvalidResponseError
	if errors.As(err, &invalidResponseError) {
		r.config.Logger.Debugf("%+v\n%s", invalidResponseError.PrivateErr, invalidResponseError.Response)
	}

	// Redirects the user to a default error URI defined in the SAML configuration.
	http.Redirect(writer, request, r.config.DefaultErrorURI, http.StatusSeeOther)
}

// ServiceProvider returns a pointer to the SAML service provider configured in the middleware.
func (r *SAML) ServiceProvider() *saml.ServiceProvider {
	return &r.middleware.ServiceProvider
}

// CreateSession is called when we have received a valid SAML assertion and
// should create a new session and modify the http response accordingly, e.g. by
// setting a cookie.
func (r *SAML) CreateSession(ctx *gin.Context, assertion *saml.Assertion) error {
	return r.middleware.Session.CreateSession(ctx.Writer, ctx.Request, assertion)
}

// DeleteSession is called to modify the response such that it removed the current
// session, e.g. by deleting a cookie.
func (r *SAML) DeleteSession(ctx *gin.Context) error {
	return r.middleware.Session.DeleteSession(ctx.Writer, ctx.Request)
}

// GetSession returns the current Session associated with the request, or
// ErrNoSession if there is no valid session.
func (r *SAML) GetSession(ctx *gin.Context) (samlsp.Session, error) {
	return r.middleware.Session.GetSession(ctx.Request)
}

// parseCookieSameSite parses the SameSite value from a string and returns the corresponding http.SameSite constant.
func parseCookieSameSite(text string) http.SameSite {
	switch strings.ToLower(text) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}
