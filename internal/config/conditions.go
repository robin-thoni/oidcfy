package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type ConditionConfig struct {
	And          *ConditionAndConfig          `yaml:"and"`
	Or           *ConditionOrConfig           `yaml:"or"`
	Not          *ConditionNotConfig          `yaml:"not"`
	Redirect     *ConditionRedirectConfig     `yaml:"redirect"`
	Unauthorized *ConditionUnauthorizedConfig `yaml:"unauthorized"`
	Host         *ConditionHostConfig         `yaml:"host"`
	Path         *ConditionPathConfig         `yaml:"path"`
	Claim        *ConditionClaimConfig        `yaml:"claim"`
}

func (cond *ConditionConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!str" {
		return yaml.Unmarshal(([]byte)(fmt.Sprintf("%s: {}", value.Value)), cond)
	}
	type rawType ConditionConfig
	return value.Decode((*rawType)(cond))
}

type ConditionAndConfig struct {
	Conditions []ConditionConfig `yaml:"conditions"`
}

func (cond *ConditionAndConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!seq" {
		return value.Decode(&cond.Conditions)
	}
	type rawType ConditionAndConfig
	return value.Decode((*rawType)(cond))
}

type ConditionOrConfig struct {
	Conditions []ConditionConfig `yaml:"conditions"`
}

func (cond *ConditionOrConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!seq" {
		return value.Decode(&cond.Conditions)
	}
	type rawType ConditionOrConfig
	return value.Decode((*rawType)(cond))
}

type ConditionNotConfig struct {
	Condition ConditionConfig `yaml:"condition"`
}

type ConditionRedirectConfig struct {
}

type ConditionUnauthorizedConfig struct {
}

type ConditionHostConfig struct {
	HostTpl string `yaml:"host"`
}

func (cond *ConditionHostConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!str" {
		cond.HostTpl = value.Value
		return nil
	}
	type rawType ConditionHostConfig
	return value.Decode((*rawType)(cond))
}

type ConditionPathConfig struct {
	PathTpl string `yaml:"path"`
}

func (cond *ConditionPathConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!str" {
		cond.PathTpl = value.Value
		return nil
	}
	type rawType ConditionPathConfig
	return value.Decode((*rawType)(cond))
}

type ConditionClaimConfig struct {
	ClaimTpl string `yaml:"claim"`
}

func (cond *ConditionClaimConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!str" {
		cond.ClaimTpl = value.Value
		return nil
	}
	type rawType ConditionClaimConfig
	return value.Decode((*rawType)(cond))
}
