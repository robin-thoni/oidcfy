package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNotSuccessSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
not:
  condition:
    true
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.Not)

	assert.NotNil(t, condConfig.Not.Condition)
	assert.NotNil(t, condConfig.Not.Condition.True)
}

func TestNotSuccessSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
not:
  condition:
    false
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.Not)

	assert.NotNil(t, condConfig.Not.Condition)
	assert.NotNil(t, condConfig.Not.Condition.False)
}

func TestNotFailSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
not2
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.Not)
}

func TestNotFailSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
not2: {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.Not)
}

func TestNotFailSimple3(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
not: {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.Not)
	assert.Nil(t, condConfig.Not.Condition)
}
