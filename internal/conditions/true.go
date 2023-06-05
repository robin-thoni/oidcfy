package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type True struct {
	Config *config.ConditionTrueConfig
}

func (cond *True) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	return true, nil
}

func (cond *True) fromConfig(condConfig *config.ConditionTrueConfig, ctx *conditionContext) []error {
	cond.Config = condConfig
	return nil
}
