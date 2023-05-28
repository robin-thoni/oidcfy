package profiles

import (
	"net/http"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/mutators"
)

type MutatorProfile struct {
	Config   *config.MutatorProfileConfig
	UsedBy   []Rule
	Mutators []interfaces.Mutator
}

func (mut *MutatorProfile) GetConfig() *config.MutatorProfileConfig {
	return mut.Config
}

func (mut *MutatorProfile) FromConfig(profileConfig *config.MutatorProfileConfig, name string) []error {
	var errs []error

	mut.Config = profileConfig
	mut.Mutators = make([]interfaces.Mutator, 0, len(profileConfig.Mutators))
	for _, config := range profileConfig.Mutators {
		mut1, errs1 := mutators.BuildFromConfig(config)
		if errs1 != nil {
			errs = append(errs, errs1...)
		}
		mut.Mutators = append(mut.Mutators, mut1)
	}

	return errs
}

func (mut *MutatorProfile) IsValid() bool {
	return mut.Mutators != nil // TODO check all list
}

func (mut *MutatorProfile) Mutate(rw http.ResponseWriter, ctx interfaces.MutatorContext) error {
	for _, mut := range mut.Mutators {
		err := mut.Mutate(rw, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
