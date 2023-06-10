package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestOrSuccessSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
or:
  conditions:
    - true
    - false
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.Or)

	assert.NotNil(t, condConfig.Or.Conditions)
	assert.Equal(t, len(condConfig.Or.Conditions), 2)
	assert.NotNil(t, condConfig.Or.Conditions[0].True)
	assert.NotNil(t, condConfig.Or.Conditions[1].False)
}

func TestOrSuccessSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
or:
  - true
  - false
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.Or)

	assert.NotNil(t, condConfig.Or.Conditions)
	assert.Equal(t, len(condConfig.Or.Conditions), 2)
	assert.NotNil(t, condConfig.Or.Conditions[0].True)
	assert.NotNil(t, condConfig.Or.Conditions[1].False)
}

func TestOrFailSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
or2
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.Or)
}

func TestOrFailSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
or2: {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.Or)
}

func TestOrFailSimple3(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
or: {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.Or)
	assert.Nil(t, condConfig.Or.Conditions)
}
