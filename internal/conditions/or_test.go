package conditions

import (
	"testing"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestOrSuccessSimple1(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Or{}
	errs := cond.fromConfig(&config.ConditionOrConfig{
		Conditions: []config.ConditionConfig{},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestOrSuccessSimple2(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Or{}
	errs := cond.fromConfig(&config.ConditionOrConfig{
		Conditions: []config.ConditionConfig{{
			True: &config.ConditionTrueConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestOrSuccessSimple3(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Or{}
	errs := cond.fromConfig(&config.ConditionOrConfig{
		Conditions: []config.ConditionConfig{{
			False: &config.ConditionFalseConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestOrSuccessSimple4(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Or{}
	errs := cond.fromConfig(&config.ConditionOrConfig{
		Conditions: []config.ConditionConfig{{
			False: &config.ConditionFalseConfig{},
		}, {
			True: &config.ConditionTrueConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestOrSuccessSimple5(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := Or{}
	errs := cond.fromConfig(&config.ConditionOrConfig{
		Conditions: []config.ConditionConfig{{
			True: &config.ConditionTrueConfig{},
		}, {
			False: &config.ConditionFalseConfig{},
		}},
	}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}
