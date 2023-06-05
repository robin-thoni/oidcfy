package conditions

import (
	"testing"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestFalseSuccessSimple1(t *testing.T) {
	condConfigCtx := conditionContext{
		Path: "<root>",
	}
	cond := False{}
	errs := cond.fromConfig(&config.ConditionFalseConfig{}, &condConfigCtx)
	assert.Empty(t, errs)

	res, err := cond.Evaluate(nil)
	assert.NoError(t, err)
	assert.False(t, res)
}
