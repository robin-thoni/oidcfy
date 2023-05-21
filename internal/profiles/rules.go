package profiles

import (
	"fmt"
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
)

type Rule struct {
	Config                    *config.RuleConfig
	MatchProfileName          *template.Template
	AuthenticationProfileName *template.Template
	AuthorizationProfileName  *template.Template
}

func (rule *Rule) GetConfig() *config.RuleConfig {
	return rule.Config
}

func (rule *Rule) GetMatchProfileName() *template.Template {
	return rule.MatchProfileName
}

func (rule *Rule) GetAuthenticationProfileName() *template.Template {
	return rule.AuthenticationProfileName
}

func (rule *Rule) GetAuthorizationProfileName() *template.Template {
	return rule.AuthorizationProfileName
}

func (rule *Rule) FromConfig(ruleConfig *config.RuleConfig, index int) []error {
	var errs []error
	var err error

	rule.Config = ruleConfig

	rule.MatchProfileName, err = template.New(fmt.Sprintf("RuleConfig.%d.MatchProfileTmpl", index)).Parse(ruleConfig.MatchProfileTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.AuthenticationProfileName, err = template.New(fmt.Sprintf("RuleConfig.%d.OidcProfileTmpl", index)).Parse(ruleConfig.OidcProfileTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.AuthorizationProfileName, err = template.New(fmt.Sprintf("RuleConfig.%d.AuthorizationProfileTmpl", index)).Parse(ruleConfig.AuthorizationProfileTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (rule *Rule) IsValid() bool {
	return rule.MatchProfileName != nil &&
		rule.AuthenticationProfileName != nil &&
		rule.AuthorizationProfileName != nil
}
