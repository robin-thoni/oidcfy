package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type And struct {
	Config     *config.ConditionAndConfig
	Conditions []interfaces.Condition
}

func (cond *And) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	for _, subCond := range cond.Conditions {
		result, err := subCond.Evaluate(ctx)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

func (cond *And) fromConfig(condConfig *config.ConditionAndConfig, ctx *conditionContext) []error {
	errs := make([]error, 0, 0)

	cond.Config = condConfig
	cond.Conditions, errs = buildFromConfigs(condConfig.Conditions, ctx)

	if len(errs) > 0 {
		return errs
	}
	return nil
}
