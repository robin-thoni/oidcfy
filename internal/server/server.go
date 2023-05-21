package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/interfaces"
	"github.com/robin-thoni/oidcfy/internal/profiles"
)

type AuthContext struct {
	RawRequest            *http.Request
	RawResponse           http.ResponseWriter
	Rule                  *profiles.Rule
	MatchProfile          *profiles.MatchProfile
	AuthenticationProfile *profiles.AuthenticationProfile
	AuthorizationProfile  *profiles.AuthorizationProfile
}

func (ctx *AuthContext) GetRawRequest() *http.Request {
	return ctx.RawRequest
}

func (ctx *AuthContext) GetRawResponse() http.ResponseWriter {
	return ctx.RawResponse
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
	RootConfig *config.RootConfig
	Profiles   *profiles.Profiles
	OidcfyMux  *http.ServeMux
}

func NewServer(rootConfig *config.RootConfig, profiles *profiles.Profiles) *Server {

	server := &Server{
		RootConfig: rootConfig,
		Profiles:   profiles,
		OidcfyMux:  http.NewServeMux(),
	}
	server.OidcfyMux.HandleFunc(rootConfig.Http.AuthForwardPath, server.authForward)
	server.OidcfyMux.HandleFunc(rootConfig.Http.AuthCallbackPath, server.authCallback)
	server.OidcfyMux.HandleFunc(rootConfig.Http.AuthLogoutPath, server.authLogout)

	return server
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	server.OidcfyMux.ServeHTTP(rw, r)
}

func (server *Server) authCallback(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	rw.Write(([]byte)("authCallback"))
}

func (server *Server) authLogout(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
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
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return // TODO
	}

	authCtx := AuthContext{
		RawRequest:  r,
		RawResponse: rw,
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
		err = authCtx.AuthenticationProfile.Authenticate(&authCtx)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
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
