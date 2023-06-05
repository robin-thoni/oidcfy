package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestTrueSuccessSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"true": {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.True)
}

func TestTrueSuccessSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"true"
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.True)
}

func TestTrueSuccessSimple3(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
true
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.True)
}

func TestTrueFailSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"true2": {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.True)
}
func TestTrueFailSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"true2"
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.True)
}
