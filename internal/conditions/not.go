package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type Not struct {
	Config    *config.ConditionNotConfig
	Condition interfaces.Condition
}

func (cond *Not) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	result, err := cond.Condition.Evaluate(ctx)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (cond *Not) fromConfig(condConfig *config.ConditionNotConfig, ctx *conditionContext) []error {
	var errs []error

	cond.Config = condConfig

	ctx1 := conditionContext{}
	ctx1.Path = "not"
	ctx1.Parent = ctx
	cond.Condition, errs = buildFromConfig(condConfig.Condition, &ctx1)

	if len(errs) > 0 {
		return errs
	}
	return nil
}
