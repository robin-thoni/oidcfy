package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type Redirect struct {
	Config *config.ConditionRedirectConfig
}

func (cond *Redirect) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	ctx.GetAuthContext().GetExtra().Oidcfy.AuthAction = interfaces.AuthActionRedirect
	return true, nil
}

func (cond *Redirect) fromConfig(condConfig *config.ConditionRedirectConfig, ctx *conditionContext) []error {
	cond.Config = condConfig
	return nil
}
