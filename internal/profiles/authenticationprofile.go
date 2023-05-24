package profiles

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"

	"github.com/coreos/go-oidc"
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
	"golang.org/x/oauth2"
)

type AuthenticationProfile struct {
	Config *config.OidcProfileConfig
	UsedBy []Rule

	DiscoveryUrl *template.Template
	ClientId     *template.Template
	ClientSecret *template.Template
	Scopes       *template.Template
}

func (rule *AuthenticationProfile) GetConfig() *config.OidcProfileConfig {
	return rule.Config
}

func (rule *AuthenticationProfile) FromConfig(profileConfig *config.OidcProfileConfig, name string) []error {
	var errs []error
	var err error

	rule.Config = profileConfig

	rule.DiscoveryUrl, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.OidcDiscoveryUrlTmpl", name)).Parse(profileConfig.Oidc.DiscoveryUrlTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.ClientId, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.OidcClientIdTmpl", name)).Parse(profileConfig.Oidc.ClientIdTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.ClientSecret, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.OidcSecretTmpl", name)).Parse(profileConfig.Oidc.ClientSecretTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.Scopes, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.ScopesTmpl", name)).Parse(profileConfig.Oidc.ScopesTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (rule *AuthenticationProfile) IsValid() bool {
	return rule.DiscoveryUrl != nil &&
		rule.ClientId != nil &&
		rule.ClientSecret != nil &&
		rule.Scopes != nil
}

func (rule *AuthenticationProfile) CheckAuthentication(ctx interfaces.AuthContext) (bool, error) {
	return false, nil
}

func (rule *AuthenticationProfile) makeOAuth2Context(ctx interfaces.AuthContext) (*oauth2.Config, *oidc.Provider, error) {
	oidcUrl, err := utils.RenderTemplate(rule.DiscoveryUrl, ctx)
	if err != nil {
		return nil, nil, err
	}
	oidcClientId, err := utils.RenderTemplate(rule.ClientId, ctx)
	if err != nil {
		return nil, nil, err
	}
	oidcClientSecret, err := utils.RenderTemplate(rule.ClientSecret, ctx)
	if err != nil {
		return nil, nil, err
	}
	oidcScopesStr, err := utils.RenderTemplate(rule.Scopes, ctx)
	if err != nil {
		return nil, nil, err
	}
	oidcScopes := strings.Split(oidcScopesStr, " ")
	provider, err := oidc.NewProvider(context.TODO(), oidcUrl)
	if err != nil {
		return nil, nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     oidcClientId,
		ClientSecret: oidcClientSecret,
		RedirectURL:  "http://127.0.0.1:8080/oidcfy/auth/callback", // TODO
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcScopes,
	}

	return &oauth2Config, provider, nil
}

func (rule *AuthenticationProfile) Authenticate(ctx interfaces.AuthContext) error {
	oauth2Config, _, err := rule.makeOAuth2Context(ctx)
	if err != nil {
		return err
	}
	state := uuid.New().String() // TODO set cookie
	nonce := uuid.New().String() // TODO set cookie
	http.Redirect(ctx.GetRawResponse(), ctx.GetRawRequest(), oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)

	ctx.GetGlobalCache().AuthCallback.Add(state, 10*time.Minute, ctx) // TODO add variable for expiration

	return nil
}

func (rule *AuthenticationProfile) AuthenticateCallback(ctx interfaces.AuthContext) error {

	// TODO verify state

	oauth2Config, provider, err := rule.makeOAuth2Context(ctx)
	if err != nil {
		return err
	}

	code := ctx.GetRawRequest().URL.Query().Get("code")
	oauth2Token, err := oauth2Config.Exchange(context.TODO(), code)
	if err != nil {
		return err
	}

	var ok bool
	ctx.GetExtra().Oidcfy.IdTokenRaw, ok = oauth2Token.Extra("id_token").(string)
	if !ok {
		return err
	}
	ctx.GetExtra().Oidcfy.AccessTokenRaw, ok = oauth2Token.Extra("access_token").(string)
	if !ok {
		return err
	}
	ctx.GetExtra().Oidcfy.RefreshTokenRaw, ok = oauth2Token.Extra("refresh_token").(string)
	if !ok {
		return err
	}

	var verifier = provider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})
	ctx.GetExtra().Oidcfy.IdToken, err = verifier.Verify(context.TODO(), ctx.GetExtra().Oidcfy.IdTokenRaw)
	if err != nil {
		return err
	}
	// TODO Verify idtoken nonce

	err = ctx.GetExtra().Oidcfy.IdToken.VerifyAccessToken(ctx.GetExtra().Oidcfy.AccessTokenRaw)
	if err != nil {
		return err
	}

	// TODO set cookie

	return nil
}
