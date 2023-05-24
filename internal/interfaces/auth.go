package interfaces

import (
	"net/http"
	"text/template"

	"github.com/coreos/go-oidc"
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
	CheckAuthentication(ctx AuthContext) (bool, error)
	Authenticate(ctx AuthContext) error
}

type AuthContextAuthorization interface {
	GetConfig() *config.AuthorizationProfileConfig
}

type AuthContextExtra struct {
	Oidcfy struct {
		AuthAction      string
		IdTokenRaw      string
		IdToken         *oidc.IDToken
		AccessTokenRaw  string
		RefreshTokenRaw string
	}
}

type AuthContext interface {
	GetRawRequest() *http.Request
	GetRawResponse() http.ResponseWriter
	GetGlobalCache() *AuthContextGlobalCache
	GetExtra() *AuthContextExtra
	GetAuthContextRule() AuthContextRule
	GetAuthContextMatch() AuthContextMatch
	GetAuthContextAuthentication() AuthContextAuthentication
	GetAuthContextAuthorization() AuthContextAuthorization
}
