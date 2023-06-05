package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFalseSuccessSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"false": {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.False)
}

func TestFalseSuccessSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"false"
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.False)
}

func TestFalseSuccessSimple3(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
false
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.NotNil(t, condConfig.False)
}

func TestFalseFailSimple1(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"false2": {}
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.False)
}
func TestFalseFailSimple2(t *testing.T) {
	condConfig := ConditionConfig{}
	condConfigStr := `
"false2"
    `
	err := yaml.Unmarshal(([]byte)(condConfigStr), &condConfig)
	assert.NoError(t, err)
	assert.Nil(t, condConfig.False)
}
