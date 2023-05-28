package profiles

import (
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/muesli/cache2go"

	"github.com/coreos/go-oidc"
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/utils"
	"golang.org/x/oauth2"
)

const (
	COOKIE_ID_TOKEN     = "oidcfy.idToken"
	COOKIE_ACCESS_TOKEN = "oidcfy.accessToken"
)

type oauthContextCache struct {
	Provider *oidc.Provider
}

type AuthenticationProfile struct {
	Config *config.OidcProfileConfig
	UsedBy []Rule

	oauthContextCache *cache2go.CacheTable
	oauthVerifyCache  *cache2go.CacheTable

	DiscoveryUrl *template.Template
	ClientId     *template.Template
	ClientSecret *template.Template
	Scopes       *template.Template

	LoginTimeout *template.Template

	CookieDomain *template.Template
	CookiePath   *template.Template
	CookieSecure *template.Template
}

func (rule *AuthenticationProfile) GetConfig() *config.OidcProfileConfig {
	return rule.Config
}

func (rule *AuthenticationProfile) FromConfig(profileConfig *config.OidcProfileConfig, name string) []error {
	var errs []error
	var err error

	rule.Config = profileConfig
	rule.oauthContextCache = cache2go.Cache("AuthenticationProfile.oauthContextCache")
	rule.oauthVerifyCache = cache2go.Cache("AuthenticationProfile.oauthVerifyCache")

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

	rule.LoginTimeout, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.LoginTimeoutTmpl", name)).Parse(profileConfig.LoginTimeoutTmpl)
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

	rule.CookieSecure, err = template.New(fmt.Sprintf("OidcProfileConfig.%s.SecureTmpl", name)).Parse(profileConfig.Cookie.SecureTmpl)
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

func parseJwt(token string) map[string]interface{} {
	payload, err := base64.RawURLEncoding.DecodeString(strings.Split(token, ".")[1])
	_ = err
	var anyJson map[string]interface{}
	json.Unmarshal(payload, &anyJson)
	return anyJson
}

func (rule *AuthenticationProfile) CheckAuthentication(rw http.ResponseWriter, r *http.Request, ctx interfaces.AuthContext) (bool, error) {

	idTokenRaw, err := r.Cookie(COOKIE_ID_TOKEN)
	if err != nil {
		return false, nil
	}
	accessTokenRaw, err := r.Cookie(COOKIE_ACCESS_TOKEN)
	if err != nil {
		return false, nil
	}
	idTokenHash := hex.EncodeToString(crypto.SHA1.New().Sum(([]byte)(idTokenRaw.Value)))
	verifyCacheItem, err := rule.oauthVerifyCache.Value(idTokenHash)
	if verifyCacheItem != nil && err == nil {
		ctx.GetExtra().Oidcfy.IdTokenRaw = idTokenRaw.Value
		ctx.GetExtra().Oidcfy.IdToken = parseJwt(idTokenRaw.Value)
		ctx.GetExtra().Oidcfy.AccessTokenRaw = accessTokenRaw.Value

		return verifyCacheItem.Data().(bool), nil
	}

	oauth2Config, provider, err := rule.makeOAuth2Context(ctx)
	if err != nil {
		return false, err
	}

	var verifier = provider.Verifier(&oidc.Config{ClientID: oauth2Config.ClientID})
	idToken, err := verifier.Verify(ctx.GetContext(), idTokenRaw.Value)
	if err != nil {
		// TODO remove cookie
		return false, nil
	}
	ctx.GetExtra().Oidcfy.IdTokenRaw = idTokenRaw.Value
	ctx.GetExtra().Oidcfy.IdToken = parseJwt(idTokenRaw.Value)
	ctx.GetExtra().Oidcfy.AccessTokenRaw = accessTokenRaw.Value
	rule.oauthVerifyCache.Add(idTokenHash, idToken.Expiry.Sub(time.Now()), true) // TODO add configuration for token caching duration, to allow shorter detection of revoked tokens
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

	var provider *oidc.Provider
	providerItem, err := rule.oauthContextCache.Value(oidcUrl)
	if err == nil && providerItem != nil {
		provider = providerItem.Data().(*oauthContextCache).Provider
	}

	if provider == nil {
		provider, err = oidc.NewProvider(ctx.GetContext(), oidcUrl)
		if err != nil {
			return nil, nil, err
		}
		rule.oauthContextCache.Add(oidcUrl, 1*time.Hour, &oauthContextCache{
			Provider: provider,
		})
	}

	oauth2Config := oauth2.Config{
		ClientID:     oidcClientId,
		ClientSecret: oidcClientSecret,
		RedirectURL:  fmt.Sprintf("%s/oidcfy/auth/callback", ctx.GetRootConfig().Http.BaseUrl),
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcScopes,
	}

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
	loginTimeout, err := utils.RenderTemplateInt(rule.LoginTimeout, ctx)
	if err != nil {
		return err
	}
	sessionExpiry := time.Duration(loginTimeout) * time.Minute
	cookiePath := baseUrl.Path + "/oidcfy/auth/callback"
	cookieDomain := strings.Split(baseUrl.Host, ":")[0]
	cookieExpire := time.Now().Local().Add(sessionExpiry)

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

	ctx.GetGlobalCache().AuthCallback.Add(state, sessionExpiry, ctx)

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
	oauth2Token, err := oauth2Config.Exchange(ctx.GetContext(), code)
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
	idToken, err := verifier.Verify(ctx.GetContext(), ctx.GetExtra().Oidcfy.IdTokenRaw)
	if err != nil {
		return err
	}

	if idToken.Nonce != nonce.Value {
		return errors.New("nonces do not match")
	}

	err = idToken.VerifyAccessToken(ctx.GetExtra().Oidcfy.AccessTokenRaw)
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
	cookieSecure, err := utils.RenderTemplateBool(rule.CookieSecure, ctx)
	if err != nil {
		return err
	}
	http.SetCookie(rw, &http.Cookie{
		Name:     COOKIE_ID_TOKEN,
		Value:    ctx.GetExtra().Oidcfy.IdTokenRaw,
		Domain:   cookieDomain,
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   cookieSecure,
		Expires:  idToken.Expiry,
	})
	http.SetCookie(rw, &http.Cookie{
		Name:     COOKIE_ACCESS_TOKEN,
		Value:    ctx.GetExtra().Oidcfy.AccessTokenRaw,
		Domain:   cookieDomain,
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   cookieSecure,
		Expires:  idToken.Expiry,
	})

	return nil
}
