package conditions

import (
	"testing"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAndSuccessSimple1(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := And{}
	errs := cond.fromConfig(&config.ConditionAndConfig{
		Conditions: []*config.ConditionConfig{},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestAndSuccessSimple2(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := And{}
	errs := cond.fromConfig(&config.ConditionAndConfig{
		Conditions: []*config.ConditionConfig{{
			True: &config.ConditionTrueConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestAndSuccessSimple3(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := And{}
	errs := cond.fromConfig(&config.ConditionAndConfig{
		Conditions: []*config.ConditionConfig{{
			False: &config.ConditionFalseConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestAndSuccessSimple4(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := And{}
	errs := cond.fromConfig(&config.ConditionAndConfig{
		Conditions: []*config.ConditionConfig{{
			False: &config.ConditionFalseConfig{},
		}, {
			True: &config.ConditionTrueConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestAndSuccessSimple5(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := And{}
	errs := cond.fromConfig(&config.ConditionAndConfig{
		Conditions: []*config.ConditionConfig{{
			True: &config.ConditionTrueConfig{},
		}, {
			False: &config.ConditionFalseConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}
