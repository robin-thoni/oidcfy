package profiles

import (
	"github.com/robin-thoni/oidcfy/internal/conditions"
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type MatchProfile struct {
	Config    *config.MatchProfileConfig
	UsedBy    []Rule
	Condition interfaces.Condition
}

func (rule *MatchProfile) GetConfig() *config.MatchProfileConfig {
	return rule.Config
}

func (rule *MatchProfile) FromConfig(profileConfig *config.MatchProfileConfig, name string) []error {
	var errs []error

	rule.Config = profileConfig
	rule.Condition, errs = conditions.BuildFromConfig(profileConfig.Condition)

	return errs
}

func (rule *MatchProfile) IsValid() bool {
	return rule.Condition != nil
}

func (rule *MatchProfile) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	return rule.Condition.Evaluate(ctx)
}
