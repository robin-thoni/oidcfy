package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/muesli/cache2go"
	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/profiles"
)

type AuthContext struct {
	RootConfig            *config.RootConfig
	RawRequest            *http.Request
	RawResponse           http.ResponseWriter
	GlobalCache           *interfaces.AuthContextGlobalCache
	Extra                 interfaces.AuthContextExtra
	Rule                  *profiles.Rule
	MatchProfile          *profiles.MatchProfile
	AuthenticationProfile *profiles.AuthenticationProfile
	AuthorizationProfile  *profiles.AuthorizationProfile
}

func (ctx *AuthContext) GetRootConfig() *config.RootConfig {
	return ctx.RootConfig
}

func (ctx *AuthContext) GetRawRequest() *http.Request {
	return ctx.RawRequest
}

func (ctx *AuthContext) GetRawResponse() http.ResponseWriter {
	return ctx.RawResponse
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
	server.OidcfyMux.HandleFunc("/oidcfy/auth/callback", server.authCallback)
	server.OidcfyMux.HandleFunc("/oidcfy/auth/logout", server.authLogout)

	return server
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	server.OidcfyMux.ServeHTTP(rw, r)
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
	ctx.RawRequest = r
	ctx.RawResponse = rw
	err = ctx.AuthenticationProfile.AuthenticateCallback(ctx)
	if err != nil {
		http.Error(rw, "bad token", http.StatusBadRequest) // TODO
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func (server *Server) authLogout(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
	rw.Write(([]byte)("authLogout"))
}

func (server *Server) authForward(rw http.ResponseWriter, r *http.Request) {

	if r.Header.Get("X-Forwarded-Method") != "" &&
		r.Header.Get("X-Forwarded-Host") != "" &&
		r.Header.Get("X-Forwarded-Uri") != "" {

		r.Method = r.Header.Get("X-Forwarded-Method")
		r.Host = r.Header.Get("X-Forwarded-Host")
		if _, ok := r.Header["X-Forwarded-Uri"]; ok {
			r.URL, _ = url.Parse(r.Header.Get("X-Forwarded-Uri"))
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return // TODO
		}
	} else if r.URL.Query().Get("X-Forwarded-Method") != "" &&
		r.URL.Query().Get("X-Forwarded-Host") != "" &&
		r.URL.Query().Get("X-Forwarded-Uri") != "" {

		r.Method = r.URL.Query().Get("X-Forwarded-Method")
		r.Host = r.URL.Query().Get("X-Forwarded-Host")
		if _, ok := r.URL.Query()["X-Forwarded-Uri"]; ok {
			r.URL, _ = url.Parse(r.URL.Query().Get("X-Forwarded-Uri"))
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return // TODO
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return // TODO
	}

	authCtx := AuthContext{
		RootConfig:  server.RootConfig,
		RawRequest:  r,
		RawResponse: rw,
		GlobalCache: server.GlobalCache,
	}

	conditionCtx := ConditionContext{
		AuthContext: &authCtx,
	}

	var err error

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

	res, err := authCtx.AuthenticationProfile.CheckAuthentication(&authCtx)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !res {
		if authCtx.GetExtra().Oidcfy.AuthAction == interfaces.AuthActionRedirect {
			err = authCtx.AuthenticationProfile.Authenticate(&authCtx)
			if err != nil {
				log.Println(err)
				rw.WriteHeader(http.StatusInternalServerError)
			}
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
	if res {
		rw.WriteHeader(http.StatusOK)
		return
	}

	log.Println("Unauthorized by profile")
	rw.WriteHeader(http.StatusUnauthorized)
}
