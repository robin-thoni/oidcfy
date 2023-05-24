package profiles

import (
	"github.com/robin-thoni/oidcfy/internal/conditions"
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type AuthorizationProfile struct {
	Config    *config.AuthorizationProfileConfig
	UsedBy    []Rule
	Condition interfaces.Condition
}

func (rule *AuthorizationProfile) GetConfig() *config.AuthorizationProfileConfig {
	return rule.Config
}

func (rule *AuthorizationProfile) FromConfig(profileConfig *config.AuthorizationProfileConfig, name string) []error {
	var errs []error

	rule.Config = profileConfig
	rule.Condition, errs = conditions.BuildFromConfig(profileConfig.Condition)

	return errs
}

func (rule *AuthorizationProfile) IsValid() bool {
	return rule.Condition != nil
}

func (rule *AuthorizationProfile) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	return rule.Condition.Evaluate(ctx)
}
