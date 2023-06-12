package conditions

import (
	"testing"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNotSuccessSimple1(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Not{}
	errs := cond.fromConfig(&config.ConditionNotConfig{
		Condition: &config.ConditionConfig{
			False: &config.ConditionFalseConfig{},
		},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestNotSuccessSimple2(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Not{}
	errs := cond.fromConfig(&config.ConditionNotConfig{
		Condition: &config.ConditionConfig{
			True: &config.ConditionTrueConfig{},
		},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestNotConfigError1(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Not{}
	errs := cond.fromConfig(&config.ConditionNotConfig{}, &condConfigCtx)
	assert.NotEmpty(t, errs)
}

func TestNotConfigError2(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Not{}
	errs := cond.fromConfig(&config.ConditionNotConfig{
		Condition: &config.ConditionConfig{},
	}, &condConfigCtx)
	assert.NotEmpty(t, errs)
}
