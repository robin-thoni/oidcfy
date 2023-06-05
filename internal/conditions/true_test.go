package conditions

import (
	"testing"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestTrueSuccessSimple1(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := True{}
	errs := cond.fromConfig(&config.ConditionTrueConfig{}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.True(t, res)
}
