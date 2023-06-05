package conditions

import (
	"errors"
	"fmt"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
)

type conditionContext struct {
	Parent *conditionContext
	Path   string
}

func (ctx *conditionContext) fullPath() string {
	if ctx.Parent != nil {
		return fmt.Sprintf("%s.%s", ctx.Parent.fullPath(), ctx.Path)
	}
	return ctx.Path
}

func BuildFromConfig(cond config.ConditionConfig) (interfaces.Condition, []error) {

	ctx := conditionContext{}
	ctx.Path = "<root>"
	return buildFromConfig(&cond, &ctx)
}

// func GetTypeForConfig[T interface{}](x T) (reflect.Type, error) {
// 	t := reflect.TypeOf(x)
// 	if t == reflect.TypeOf((*config.ConditionAndConfig)(nil)).Elem() {
// 		return reflect.TypeOf((*And)(nil)).Elem(), nil
// 	}
// 	return nil, errors.New(fmt.Sprintf("Unknown condition config type: %s", t))
// }

func buildFromConfigs(condConfigs []config.ConditionConfig, ctx *conditionContext) ([]interfaces.Condition, []error) {
	errs := make([]error, 0, 0)
	conds := make([]interfaces.Condition, 0, len(condConfigs))
	for i, cond := range condConfigs {
		ctx1 := conditionContext{}
		ctx1.Path = fmt.Sprintf("%d", i)
		ctx1.Parent = ctx
		cond1, errs1 := buildFromConfig(&cond, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, cond1) // append even if nil
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return conds, nil
}

func buildFromConfig(condConfig *config.ConditionConfig, ctx *conditionContext) (interfaces.Condition, []error) {

	errs := make([]error, 0, 0)
	conds := make([]interfaces.Condition, 0, 1)

	if condConfig.True != nil {
		cond1 := True{}
		errs1 := cond1.fromConfig(condConfig.True, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.And != nil {
		cond1 := And{}
		errs1 := cond1.fromConfig(condConfig.And, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Or != nil {
		cond1 := Or{}
		errs1 := cond1.fromConfig(condConfig.Or, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Not != nil {
		cond1 := Not{}
		errs1 := cond1.fromConfig(condConfig.Not, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Redirect != nil {
		cond1 := Redirect{}
		errs1 := cond1.fromConfig(condConfig.Redirect, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Unauthorized != nil {
		cond1 := Unauthorized{}
		errs1 := cond1.fromConfig(condConfig.Unauthorized, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Host != nil {
		cond1 := Host{}
		errs1 := cond1.fromConfig(condConfig.Host, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Path != nil {
		cond1 := Path{}
		errs1 := cond1.fromConfig(condConfig.Path, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}
	if condConfig.Claim != nil {
		cond1 := Claim{}
		errs1 := cond1.fromConfig(condConfig.Claim, ctx)
		if len(errs1) > 0 {
			errs = append(errs, errs1...)
		}
		conds = append(conds, &cond1)
	}

	// v := reflect.ValueOf(condConfig)
	// for i := 0; i < v.NumField(); i++ {
	// 	if !v.Field(i).IsNil() {
	// 		t, err := GetTypeForConfig(v.Field(i).Interface())
	// 		if err != nil {
	// 			errs = append(errs, err)
	// 		}
	// 		if t != nil {
	// 			cond := reflect.New(t).Elem().Interface().(interfaces.Condition)
	// 			errs1 := cond.FromConfig(condConfig.And, ctx)
	// 		}
	// 		// if foundNonNilCondition != nil {
	// 		// 	return false, errors.New("Multiple conditions defined")
	// 		// }
	// 		// foundNonNilCondition = v.Field(i).Interface().(Condition)
	// 	}
	// }

	if len(conds) == 0 {
		errs = append(errs, errors.New(fmt.Sprintf("%s: No condition defined", ctx.fullPath())))
	}
	if len(conds) > 1 {
		errs = append(errs, errors.New(fmt.Sprintf("%s: Multiple conditions defined", ctx.fullPath())))
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return conds[0], nil
}
