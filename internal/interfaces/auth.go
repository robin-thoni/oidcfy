package interfaces

import (
	"context"
	"net/http"
	"net/url"
	"text/template"

	"github.com/muesli/cache2go"
	"github.com/robin-thoni/oidcfy/internal/config"
)

const (
	AuthActionRedirect     = "redirect"
	AuthActionUnauthorized = "unauthorized"
)

type AuthContextGlobalCache struct {
	AuthCallback *cache2go.CacheTable
}

type AuthOriginalRequest struct {
	Url    url.URL
	Method string
}

type AuthContextRule interface {
	GetConfig() *config.RuleConfig
	GetMatchProfileName() *template.Template
	GetAuthenticationProfileName() *template.Template
	GetAuthorizationProfileName() *template.Template
}

type AuthContextMatch interface {
	GetConfig() *config.MatchProfileConfig
}

type AuthContextAuthentication interface {
	GetConfig() *config.OidcProfileConfig
	CheckAuthentication(rw http.ResponseWriter, r *http.Request, ctx AuthContext) (bool, error)
	Authenticate(rw http.ResponseWriter, r *http.Request, ctx AuthContext) error
}

type AuthContextAuthorization interface {
	GetConfig() *config.AuthorizationProfileConfig
}

type AuthContextExtra struct {
	Oidcfy struct {
		AuthAction      string
		IdTokenRaw      string
		IdToken         map[string]interface{}
		AccessTokenRaw  string
		AccessToken     map[string]interface{}
		RefreshTokenRaw string
	}
}

type AuthContext interface {
	GetContext() context.Context
	GetRootConfig() *config.RootConfig
	GetOriginalRequest() *AuthOriginalRequest
	GetGlobalCache() *AuthContextGlobalCache
	GetExtra() *AuthContextExtra
	GetAuthContextRule() AuthContextRule
	GetAuthContextMatch() AuthContextMatch
	GetAuthContextAuthentication() AuthContextAuthentication
	GetAuthContextAuthorization() AuthContextAuthorization
}
