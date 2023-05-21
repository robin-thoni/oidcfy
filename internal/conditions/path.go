package conditions

import (
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
)

type Path struct {
	Config  *config.ConditionPathConfig
	PathTpl *template.Template
}

func (cond *Path) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	host, err := utils.RenderTemplate(cond.PathTpl, ctx)
	if err != nil {
		return false, err
	}
	return ctx.GetAuthContext().GetRawRequest().URL.Path == host, nil
}

func (cond *Path) fromConfig(condConfig *config.ConditionPathConfig, ctx *conditionContext) []error {
	var err error

	cond.Config = condConfig
	cond.PathTpl, err = template.New(ctx.fullPath()).Parse(condConfig.PathTpl)

	if err != nil {
		return []error{err}
	}
	return nil
}
