package profiles

import (
	"errors"
	"fmt"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
)

type Profiles struct {
	MatchProfiles          map[string]*MatchProfile
	AuthenticationProfiles map[string]*AuthenticationProfile
	AuthorizationProfiles  map[string]*AuthorizationProfile
	Rules                  []*Rule
}

func (profs *Profiles) FromConfig(rootConfig *config.RootConfig) []error {
	errs := []error{}

	profs.MatchProfiles = map[string]*MatchProfile{}
	for profileName, profileConfig := range rootConfig.MatchProfiles {
		matchProfile := MatchProfile{}
		errs1 := matchProfile.FromConfig(&profileConfig, profileName)
		if len(errs) > 0 {
			errs = append(errs, errs1...)
		}
		profs.MatchProfiles[profileName] = &matchProfile // append even if errors
	}

	profs.AuthenticationProfiles = map[string]*AuthenticationProfile{}
	for profileName, profileConfig := range rootConfig.OidcProfiles {
		matchProfile := AuthenticationProfile{}
		errs1 := matchProfile.FromConfig(&profileConfig, profileName)
		if len(errs) > 0 {
			errs = append(errs, errs1...)
		}
		profs.AuthenticationProfiles[profileName] = &matchProfile // append even if errors
	}

	profs.AuthorizationProfiles = map[string]*AuthorizationProfile{}
	for profileName, profileConfig := range rootConfig.AuthorizationProfiles {
		matchProfile := AuthorizationProfile{}
		errs1 := matchProfile.FromConfig(&profileConfig, profileName)
		if len(errs) > 0 {
			errs = append(errs, errs1...)
		}
		profs.AuthorizationProfiles[profileName] = &matchProfile // append even if errors
	}

	profs.Rules = []*Rule{}
	for ruleConfigIndex, ruleConfig := range rootConfig.Rules {
		rule := Rule{}
		errs1 := rule.FromConfig(&ruleConfig, ruleConfigIndex)
		if len(errs) > 0 {
			errs = append(errs, errs1...)
		}
		profs.Rules = append(profs.Rules, &rule) // append even if errors
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (profiles *Profiles) GetMatchProfile(name string) (*MatchProfile, error) {
	matchProfile, ok := profiles.MatchProfiles[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No match profile named: %s", name))
	}
	if matchProfile == nil || !matchProfile.IsValid() {
		return nil, errors.New(fmt.Sprintf("Match profile %s is invalid", name))
	}

	return matchProfile, nil
}

func (profiles *Profiles) GetRule(ctx interfaces.ConditionContext) (*Rule, *MatchProfile, error) {
	for _, rule := range profiles.Rules {
		name, err := utils.RenderTemplate(rule.MatchProfileName, ctx.GetAuthContext())
		if err != nil {
			return nil, nil, err
		}
		matchProfile, err := profiles.GetMatchProfile(name)
		if err != nil {
			return nil, nil, err
		}
		if matchProfile != nil && matchProfile.IsValid() {
			res, err := matchProfile.Evaluate(ctx)
			if err != nil {
				return nil, nil, err
			}
			if res {
				return rule, matchProfile, nil
			}
		}
	}
	return nil, nil, nil
}

func (profiles *Profiles) GetAuthenticationProfile(ctx interfaces.AuthContext) (*AuthenticationProfile, error) {
	rule := ctx.GetAuthContextRule()
	name, err := utils.RenderTemplate(rule.GetAuthenticationProfileName(), ctx)
	if err != nil {
		return nil, err
	}

	profile, ok := profiles.AuthenticationProfiles[name]
	if !ok {
		return nil, nil //errors.New(fmt.Sprintf("No authentication profile named: %s", name))
	}
	if profile == nil || !profile.IsValid() {
		return nil, errors.New(fmt.Sprintf("Authentication profile %s is invalid", name))
	}

	return profile, nil
}

func (profiles *Profiles) GetAuthorizationProfile(ctx interfaces.AuthContext) (*AuthorizationProfile, error) {
	rule := ctx.GetAuthContextRule()
	name, err := utils.RenderTemplate(rule.GetAuthorizationProfileName(), ctx)
	if err != nil {
		return nil, err
	}

	profile, ok := profiles.AuthorizationProfiles[name]
	if !ok {
		return nil, nil //errors.New(fmt.Sprintf("No authorization profile named: %s", name))
	}
	if profile == nil || !profile.IsValid() {
		return nil, errors.New(fmt.Sprintf("Authorization profile %s is invalid", name))
	}

	return profile, nil
}
