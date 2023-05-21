package conditions

import (
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
)

type Host struct {
	Config  *config.ConditionHostConfig
	HostTpl *template.Template
}

func (cond *Host) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	host, err := utils.RenderTemplate(cond.HostTpl, ctx)
	if err != nil {
		return false, err
	}
	return ctx.GetAuthContext().GetRawRequest().Host == host, nil
}

func (cond *Host) fromConfig(condConfig *config.ConditionHostConfig, ctx *conditionContext) []error {
	var err error

	cond.Config = condConfig
	cond.HostTpl, err = template.New(ctx.fullPath()).Parse(condConfig.HostTpl)

	if err != nil {
		return []error{err}
	}
	return nil
}
