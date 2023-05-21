package conditions

import (
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type Or struct {
	Config     *config.ConditionOrConfig
	Conditions []interfaces.Condition
}

func (cond *Or) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	for _, subCond := range cond.Conditions {
		result, err := subCond.Evaluate(ctx)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func (cond *Or) fromConfig(condConfig *config.ConditionOrConfig, ctx *conditionContext) []error {
	errs := make([]error, 0, 0)

	cond.Config = condConfig
	cond.Conditions, errs = buildFromConfigs(condConfig.Conditions, ctx)

	if len(errs) > 0 {
		return errs
	}
	return nil
}
