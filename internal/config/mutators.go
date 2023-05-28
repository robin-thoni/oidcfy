package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type MutatorConfig struct {
	Header *MutatorHeaderConfig `yaml:"header"`
}

func (mut *MutatorConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!str" {
		return yaml.Unmarshal(([]byte)(fmt.Sprintf("%s: {}", value.Value)), mut)
	}
	type rawType MutatorConfig
	return value.Decode((*rawType)(mut))
}

type MutatorHeaderConfig struct {
	NameTpl  string `yaml:"name"`
	ValueTpl string `yaml:"value"`
}
