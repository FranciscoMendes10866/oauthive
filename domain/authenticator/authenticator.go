package authenticator

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
)

type Authenticator interface {
	InitializeLogin(provider string, w http.ResponseWriter, r *http.Request)
	CompleteLogin(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error)
}

type authenticator struct {
	discordProvider *discord.Provider
}

func NewAuthenticator() Authenticator {
	discordProvider := discord.New(
		os.Getenv("DISCORD_CLIENT_ID"),
		os.Getenv("DISCORD_CLIENT_SECRET"),
		os.Getenv("DISCORD_CALLBACK_URL"),
		[]string{discord.ScopeIdentify, discord.ScopeEmail}...,
	)

	goth.UseProviders([]goth.Provider{discordProvider}...)

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.MaxAge(60 * 4) // 4 minutes
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("API_ENV") == "prod"

	gothic.Store = store

	return &authenticator{
		discordProvider: discordProvider,
	}
}

func (self *authenticator) InitializeLogin(provider string, w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (self *authenticator) CompleteLogin(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error) {
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	return gothic.CompleteUserAuth(w, r)
}
