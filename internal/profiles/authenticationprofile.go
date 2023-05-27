package profiles

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

	CookieDomain *template.Template
	CookiePath   *template.Template
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

	rule.CookieDomain, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.DomainTmpl", name)).Parse(profileConfig.Cookie.DomainTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	rule.CookiePath, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.PathTmpl", name)).Parse(profileConfig.Cookie.PathTmpl)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (rule *AuthenticationProfile) IsValid() bool {
	return rule.DiscoveryUrl != nil &&
		rule.ClientId != nil &&
		rule.ClientSecret != nil &&
		rule.Scopes != nil &&
		rule.CookieDomain != nil &&
		rule.CookiePath != nil
}

func (rule *AuthenticationProfile) CheckAuthentication(rw http.ResponseWriter, r *http.Request, ctx interfaces.AuthContext) (bool, error) {

	idToken, err := r.Cookie("idToken")
	if err != nil {
		return false, nil
	}

	oauth2Config, provider, err := rule.makeOAuth2Context(ctx)
	if err != nil {
		return false, err
	}
	var verifier = provider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})
	_, err = verifier.Verify(context.TODO(), idToken.Value)
	if err != nil {
		// TODO remove cookie
		return false, nil
	}
	return true, nil
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
		RedirectURL:  fmt.Sprintf("%s/oidcfy/auth/callback", ctx.GetRootConfig().Http.BaseUrl),
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcScopes,
	}

	// TODO cache provider and oauth2Config

	return &oauth2Config, provider, nil
}

func (rule *AuthenticationProfile) Authenticate(rw http.ResponseWriter, r *http.Request, ctx interfaces.AuthContext) error {
	oauth2Config, _, err := rule.makeOAuth2Context(ctx)
	if err != nil {
		return err
	}

	baseUrl, err := url.Parse(ctx.GetRootConfig().Http.BaseUrl)
	if err != nil {
		return err
	}
	cookiePath := baseUrl.Path + "/oidcfy/auth/callback"
	cookieDomain := strings.Split(baseUrl.Host, ":")[0]
	cookieExpire := time.Now().Local().Add(10 * time.Minute) // TODO add variable for expiration

	state := uuid.New().String()
	http.SetCookie(rw, &http.Cookie{
		Name:     "state",
		Value:    state,
		Path:     cookiePath,
		Domain:   cookieDomain,
		HttpOnly: true,
		Secure:   baseUrl.Scheme == "https",
		Expires:  cookieExpire,
	})

	nonce := uuid.New().String()
	http.SetCookie(rw, &http.Cookie{
		Name:     "nonce",
		Value:    nonce,
		Path:     cookiePath,
		Domain:   cookieDomain,
		HttpOnly: true,
		Secure:   baseUrl.Scheme == "https",
		Expires:  cookieExpire,
	})

	http.Redirect(rw, r, oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)

	ctx.GetGlobalCache().AuthCallback.Add(state, 10*time.Minute, ctx) // TODO add variable for expiration

	return nil
}

func (rule *AuthenticationProfile) AuthenticateCallback(rw http.ResponseWriter, r *http.Request, ctx interfaces.AuthContext) error {

	state, err := r.Cookie("state")
	if err != nil {
		return err
	}
	nonce, err := r.Cookie("nonce")
	if err != nil {
		return err
	}

	if state.Value != r.URL.Query().Get("state") {
		return errors.New("states do not match")
	}

	oauth2Config, provider, err := rule.makeOAuth2Context(ctx)
	if err != nil {
		return err
	}

	code := r.URL.Query().Get("code")
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

	if ctx.GetExtra().Oidcfy.IdToken.Nonce != nonce.Value {
		return errors.New("nonces do not match")
	}

	err = ctx.GetExtra().Oidcfy.IdToken.VerifyAccessToken(ctx.GetExtra().Oidcfy.AccessTokenRaw)
	if err != nil {
		return err
	}

	baseUrl, err := url.Parse(ctx.GetRootConfig().Http.BaseUrl)
	if err != nil {
		return err
	}
	cookieDomain, err := utils.RenderTemplate(rule.CookieDomain, ctx)
	if err != nil {
		return err
	}
	cookiePath, err := utils.RenderTemplate(rule.CookiePath, ctx)
	if err != nil {
		return err
	}
	http.SetCookie(rw, &http.Cookie{
		Name:     "idToken",
		Value:    ctx.GetExtra().Oidcfy.IdTokenRaw,
		Domain:   cookieDomain,
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   baseUrl.Scheme == "https", // TODO must match with applications
		Expires:  ctx.GetExtra().Oidcfy.IdToken.Expiry,
	})
	http.SetCookie(rw, &http.Cookie{
		Name:     "accessToken",
		Value:    ctx.GetExtra().Oidcfy.AccessTokenRaw,
		Domain:   cookieDomain,
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   baseUrl.Scheme == "https", // TODO must match with applications
		Expires:  ctx.GetExtra().Oidcfy.IdToken.Expiry,
	})

	return nil
}
