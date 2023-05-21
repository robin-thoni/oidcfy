package profiles

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type AuthenticationProfile struct {
	Config *config.OidcProfileConfig
	UsedBy []Rule

	OidcDiscoveryUrl *template.Template
	OidcClientId     *template.Template
	OidcSecret       *template.Template
}

func (rule *AuthenticationProfile) FromConfig(profileConfig *config.OidcProfileConfig, name string) []error {
	var errs []error
	var err error

	rule.Config = profileConfig

	rule.OidcDiscoveryUrl, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.OidcDiscoveryUrlTmpl", name)).Parse(profileConfig.OidcDiscoveryUrlTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.OidcClientId, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.OidcClientIdTmpl", name)).Parse(profileConfig.OidcClientIdTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.OidcSecret, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.OidcSecretTmpl", name)).Parse(profileConfig.OidcSecretTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (rule *AuthenticationProfile) IsValid() bool {
	return rule.OidcDiscoveryUrl != nil &&
		rule.OidcClientId != nil &&
		rule.OidcSecret != nil
}

func (rule *AuthenticationProfile) CheckAuthentication(ctx interfaces.AuthContext) (bool, error) {
	return false, nil
}

func (rule *AuthenticationProfile) Authenticate(ctx interfaces.AuthContext) error {
	// oidcUrl, err := utils.RenderTemplate(rule.OidcDiscoveryUrl, ctx)// TODO use manifest auth endpoint
	// if err != nil {
	// 	return err
	// }
	// ctx.GetRawResponse().Header().Add("Location", oidcUrl)
	// ctx.GetRawResponse().WriteHeader(http.StatusTemporaryRedirect)
	ctx.GetRawResponse().WriteHeader(http.StatusNoContent)
	return nil //errors.New("Failed to authenticate")
}
