package interfaces

import (
	"net/http"
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
)

type AuthContextRule interface {
	GetConfig() *config.RuleConfig
	GetMatchProfileName() *template.Template
	GetAuthenticationProfileName() *template.Template
	GetAuthorizationProfileName() *template.Template
}

type AuthContextMatch interface {
}

type AuthContextAuthentication interface {
	CheckAuthentication(ctx AuthContext) (bool, error)
}

type AuthContextAuthorization interface {
}

type AuthContextExtra struct {
	Oidcfy struct {
		AuthAction string
	}
}

type AuthContext interface {
	GetRawRequest() *http.Request
	GetRawResponse() http.ResponseWriter
	GetExtra() *AuthContextExtra
	GetAuthContextRule() AuthContextRule
	GetAuthContextMatch() AuthContextMatch
	GetAuthContextAuthentication() AuthContextAuthentication
	GetAuthContextAuthorization() AuthContextAuthorization
}
