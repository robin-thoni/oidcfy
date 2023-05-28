package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/muesli/cache2go"
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/profiles"
)

type AuthContext struct {
	Context               context.Context
	RootConfig            *config.RootConfig
	OriginalRequest       *interfaces.AuthOriginalRequest
	GlobalCache           *interfaces.AuthContextGlobalCache
	Extra                 interfaces.AuthContextExtra
	Rule                  *profiles.Rule
	MatchProfile          *profiles.MatchProfile
	AuthenticationProfile *profiles.AuthenticationProfile
	AuthorizationProfile  *profiles.AuthorizationProfile
	MutatorProfile        *profiles.MutatorProfile
}

func (ctx *AuthContext) GetContext() context.Context {
	return ctx.Context
}

func (ctx *AuthContext) GetRootConfig() *config.RootConfig {
	return ctx.RootConfig
}

func (ctx *AuthContext) GetOriginalRequest() *interfaces.AuthOriginalRequest {
	return ctx.OriginalRequest
}

func (ctx *AuthContext) GetGlobalCache() *interfaces.AuthContextGlobalCache {
	return ctx.GlobalCache
}

func (ctx *AuthContext) GetExtra() *interfaces.AuthContextExtra {
	return &ctx.Extra
}

func (ctx *AuthContext) GetAuthContextRule() interfaces.AuthContextRule {
	return ctx.Rule
}

func (ctx *AuthContext) GetAuthContextMatch() interfaces.AuthContextMatch {
	return ctx.MatchProfile
}

func (ctx *AuthContext) GetAuthContextAuthentication() interfaces.AuthContextAuthentication {
	return ctx.AuthenticationProfile
}

func (ctx *AuthContext) GetAuthContextAuthorization() interfaces.AuthContextAuthorization {
	return ctx.AuthorizationProfile
}

func (ctx *AuthContext) GetAuthContextMutator() interfaces.AuthContextMutator {
	return ctx.MutatorProfile
}

type ConditionContextDebug struct {
}

type ConditionContext struct {
	AuthContext interfaces.AuthContext
	Debug       interfaces.ConditionContextDebug
}

func (ctx *ConditionContext) GetAuthContext() interfaces.AuthContext {
	return ctx.AuthContext
}

func (ctx *ConditionContext) GetDebug() interfaces.ConditionContextDebug {
	return ctx.Debug
}

type MutatorContextDebug struct {
}

type MutatorContext struct {
	AuthContext interfaces.AuthContext
	Debug       interfaces.MutatorContextDebug
}

func (ctx *MutatorContext) GetAuthContext() interfaces.AuthContext {
	return ctx.AuthContext
}

func (ctx *MutatorContext) GetDebug() interfaces.MutatorContextDebug {
	return ctx.Debug
}

type Server struct {
	RootConfig  *config.RootConfig
	Profiles    *profiles.Profiles
	OidcfyMux   *http.ServeMux
	GlobalCache *interfaces.AuthContextGlobalCache
}

func NewServer(rootConfig *config.RootConfig, profiles *profiles.Profiles) *Server {

	server := &Server{
		RootConfig: rootConfig,
		Profiles:   profiles,
		OidcfyMux:  http.NewServeMux(),
		GlobalCache: &interfaces.AuthContextGlobalCache{
			AuthCallback: cache2go.Cache("oidcfy"),
		},
	}
	server.OidcfyMux.HandleFunc("/oidcfy/auth/forward", server.authForward)
	server.OidcfyMux.HandleFunc("/oidcfy/auth/login", server.authLogin)
	server.OidcfyMux.HandleFunc("/oidcfy/auth/callback", server.authCallback)
	server.OidcfyMux.HandleFunc("/oidcfy/auth/logout", server.authLogout)

	return server
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Printf("%s %s", r.URL, elapsed)
	}()
	server.OidcfyMux.ServeHTTP(rw, r)
}

func buildOriginalRequest(r *http.Request, fromQuery bool) (interfaces.AuthOriginalRequest, error) {
	req := interfaces.AuthOriginalRequest{}
	if r.Header.Get("X-Forwarded-Method") != "" &&
		r.Header.Get("X-Forwarded-Proto") != "" &&
		r.Header.Get("X-Forwarded-Host") != "" &&
		r.Header.Get("X-Forwarded-Uri") != "" {

		req.Method = r.Header.Get("X-Forwarded-Method")
		uri, err := url.Parse(r.Header.Get("X-Forwarded-Uri"))
		if err != nil {
			return req, err
		}
		req.Url = *uri
		req.Url.Scheme = r.Header.Get("X-Forwarded-Proto")
		req.Url.Host = r.Header.Get("X-Forwarded-Host")
	} else if fromQuery &&
		r.URL.Query().Get("X-Forwarded-Method") != "" &&
		r.URL.Query().Get("X-Forwarded-Proto") != "" &&
		r.URL.Query().Get("X-Forwarded-Host") != "" &&
		r.URL.Query().Get("X-Forwarded-Uri") != "" {

		req.Method = r.URL.Query().Get("X-Forwarded-Method")
		uri, err := url.Parse(r.URL.Query().Get("X-Forwarded-Uri"))
		if err != nil {
			return req, err
		}
		req.Url = *uri
		req.Url.Scheme = r.URL.Query().Get("X-Forwarded-Proto")
		req.Url.Host = r.URL.Query().Get("X-Forwarded-Host")
	} else {
		return req, errors.New("Missing X-Forwarded-* parameters")
	}
	return req, nil
}

func (server *Server) authLogin(rw http.ResponseWriter, r *http.Request) {
	originalRequest, err := buildOriginalRequest(r, true)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	authCtx := AuthContext{
		Context:         context.Background(),
		RootConfig:      server.RootConfig,
		OriginalRequest: &originalRequest,
		GlobalCache:     server.GlobalCache,
	}

	conditionCtx := ConditionContext{
		AuthContext: &authCtx,
	}

	authCtx.Rule, authCtx.MatchProfile, err = server.Profiles.GetRule(&conditionCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authCtx.Rule == nil {
		log.Println("Failed to find a matching rule")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	authCtx.AuthenticationProfile, err = server.Profiles.GetAuthenticationProfile(&authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authCtx.AuthenticationProfile == nil {
		log.Println("Failed to find a matching authentication profile")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = authCtx.AuthenticationProfile.Authenticate(rw, r, &authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *Server) authCallback(rw http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(rw, "state not found", http.StatusBadRequest) // TODO
		return
	}

	cacheItem, err := server.GlobalCache.AuthCallback.Value(state)
	if err != nil {
		http.Error(rw, "session not found", http.StatusBadRequest) // TODO
		return
	}
	server.GlobalCache.AuthCallback.Delete(state)

	ctx := cacheItem.Data().(*AuthContext)
	ctx.Context = context.Background()
	err = ctx.AuthenticationProfile.AuthenticateCallback(rw, r, ctx)
	if err != nil {
		http.Error(rw, "bad token", http.StatusBadRequest) // TODO
		return
	}

	http.Redirect(rw, r, ctx.OriginalRequest.Url.String(), http.StatusFound)
}

func (server *Server) authLogout(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
	rw.Write(([]byte)("authLogout"))
}

func (server *Server) authForward(rw http.ResponseWriter, r *http.Request) {
	originalRequest, err := buildOriginalRequest(r, true)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	authCtx := AuthContext{
		Context:         context.Background(),
		RootConfig:      server.RootConfig,
		OriginalRequest: &originalRequest,
		GlobalCache:     server.GlobalCache,
	}

	conditionCtx := ConditionContext{
		AuthContext: &authCtx,
	}

	authCtx.Rule, authCtx.MatchProfile, err = server.Profiles.GetRule(&conditionCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authCtx.Rule == nil {
		log.Println("Failed to find a matching rule")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	authCtx.AuthenticationProfile, err = server.Profiles.GetAuthenticationProfile(&authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authCtx.AuthenticationProfile == nil {
		log.Println("Failed to find a matching authentication profile")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := authCtx.AuthenticationProfile.CheckAuthentication(rw, r, &authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !res {
		if authCtx.GetExtra().Oidcfy.AuthAction == interfaces.AuthActionRedirect {
			url, err := url.Parse(server.RootConfig.Http.BaseUrl + "/oidcfy/auth/login")
			if err != nil {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			query := url.Query()
			query.Add("X-Forwarded-Method", originalRequest.Method)
			query.Add("X-Forwarded-Proto", originalRequest.Url.Scheme)
			query.Add("X-Forwarded-Host", originalRequest.Url.Host)
			query.Add("X-Forwarded-Uri", originalRequest.Url.String())
			url.RawQuery = query.Encode()
			http.Redirect(rw, r, url.String(), http.StatusFound)
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	authCtx.AuthorizationProfile, err = server.Profiles.GetAuthorizationProfile(&authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authCtx.AuthorizationProfile == nil {
		log.Println("Failed to find a matching authorization profile")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err = authCtx.AuthorizationProfile.Evaluate(&conditionCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !res {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	mutatorCtx := MutatorContext{
		AuthContext: &authCtx,
	}

	authCtx.MutatorProfile, err = server.Profiles.GetMutatorProfile(&authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authCtx.MutatorProfile == nil {
		log.Println("Failed to find a matching mutator profile")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = authCtx.MutatorProfile.Mutate(rw, &mutatorCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
