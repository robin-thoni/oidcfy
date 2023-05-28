package mutators

import (
	"net/http"
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
)

type Header struct {
	Config   *config.MutatorHeaderConfig
	NameTpl  *template.Template
	ValueTpl *template.Template
}

func (cond *Header) Mutate(rw http.ResponseWriter, ctx interfaces.MutatorContext) error {
	name, err := utils.RenderTemplate(cond.NameTpl, ctx.GetAuthContext())
	if err != nil {
		return err
	}
	value, err := utils.RenderTemplate(cond.ValueTpl, ctx.GetAuthContext())
	if err != nil {
		return err
	}
	rw.Header().Add(name, value)
	return nil
}

func (cond *Header) fromConfig(condConfig *config.MutatorHeaderConfig) []error {
	var err error

	cond.Config = condConfig
	cond.NameTpl, err = template.New("Header.NameTpl").Parse(condConfig.NameTpl)
	cond.ValueTpl, err = template.New("Header.ValueTpl").Parse(condConfig.ValueTpl)

	if err != nil {
		return []error{err}
	}
	return nil
}
