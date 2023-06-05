package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type False struct {
	Config *config.ConditionFalseConfig
}

func (cond *False) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	return false, nil
}

func (cond *False) fromConfig(condConfig *config.ConditionFalseConfig, ctx *conditionContext) []error {
	cond.Config = condConfig
	return nil
}
