package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type Unauthorized struct {
	Config *config.ConditionUnauthorizedConfig
}

func (cond *Unauthorized) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	ctx.GetAuthContext().GetExtra().Oidcfy.AuthAction = interfaces.AuthActionUnauthorized
	return true, nil
}

func (cond *Unauthorized) fromConfig(condConfig *config.ConditionUnauthorizedConfig, ctx *conditionContext) []error {
	cond.Config = condConfig
	return nil
}
