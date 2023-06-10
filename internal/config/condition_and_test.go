package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestAndSuccessSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
and:
  conditions:
    - true
    - false
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.And)

	assert.NotNil(t, condConfig.And.Conditions)
	assert.Equal(t, len(condConfig.And.Conditions), 2)
	assert.NotNil(t, condConfig.And.Conditions[0].True)
	assert.NotNil(t, condConfig.And.Conditions[1].False)
}

func TestAndSuccessSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
and:
  - true
  - false
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.And)

	assert.NotNil(t, condConfig.And.Conditions)
	assert.Equal(t, len(condConfig.And.Conditions), 2)
	assert.NotNil(t, condConfig.And.Conditions[0].True)
	assert.NotNil(t, condConfig.And.Conditions[1].False)
}

func TestAndFailSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
and2
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.And)
}

func TestAndFailSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
and2: {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.And)
}

func TestAndFailSimple3(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
and: {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.And)
	assert.Nil(t, condConfig.And.Conditions)
}
