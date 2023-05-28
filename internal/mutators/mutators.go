package mutators

import (
	"errors"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

func BuildFromConfig(mutConfig config.MutatorConfig) (interfaces.Mutator, []error) {

	errs := make([]error, 0, 0)
	muts := make([]interfaces.Mutator, 0, 1)

	if mutConfig.Header != nil {
		mut1 := Header{}
		errs1 := mut1.fromConfig(mutConfig.Header)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		muts = append(muts, &mut1)
	}

	if len(muts) == 0 {
		errs = append(errs, errors.New("No mutator defined"))
	}
	if len(muts) > 1 {
		errs = append(errs, errors.New("Multiple mutators defined"))
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return muts[0], nil
}
