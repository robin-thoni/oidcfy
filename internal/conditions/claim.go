package conditions

import (
	"strings"
	"text/template"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
)

type Claim struct {
	Config   *config.ConditionClaimConfig
	ClaimTpl *template.Template
}

func (cond *Claim) Evaluate(ctx interfaces.ConditionContext) (bool, error) {
	claim, err := utils.RenderTemplate(cond.ClaimTpl, ctx.GetAuthContext())
	if err != nil {
		return false, err
	}
	claims := strings.Split(claim, ".")
	claimsCount := len(claims)

	obj := &ctx.GetAuthContext().GetExtra().Oidcfy.IdToken

	for i, c := range claims {
		subObj, ok := (*obj)[c]
		if !ok {
			break
		}
		if i < claimsCount-2 {
			tmp, ok := subObj.(map[string]interface{})
			if !ok {
				break
			}
			obj = &tmp
		} else {
			tmp, ok := subObj.([]interface{})
			if !ok {
				break
			}
			for _, item := range tmp {
				if item == claims[claimsCount-1] {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (cond *Claim) fromConfig(condConfig *config.ConditionClaimConfig, ctx *conditionContext) []error {
	var err error

	cond.Config = condConfig
	cond.ClaimTpl, err = template.New(ctx.fullPath()).Parse(condConfig.ClaimTpl)

	if err != nil {
		return []error{err}
	}
	return nil
}
